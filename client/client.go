package client

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/koron/go-mqtt/internal/waitop"
	"github.com/koron/go-mqtt/packet"
)

// Client is a MQTT client.
type Client interface {
	// Disconnect shutdown MQTT connection.
	Disconnect(force bool) error

	// Ping sends a PING packet.
	Ping() error

	// Subscribe subsribes to topics.
	Subscribe(topics []Topic) error

	// Unsubscribe unsubscribes from topics.
	Unsubscribe(topics []string) error

	// Publish publishes a message to MQTT broker.
	Publish(qos QoS, retain bool, topic string, msg []byte) error

	// Read returns a message if it was available.
	// If any messages are unavailable, this blocks until message would be
	// available when block is true, and this returns nil when block is false.
	Read(block bool) (*Message, error)
}

var (
	// ErrUnknownProtocol indicates connect adddress includes unknown protocol.
	ErrUnknownProtocol = errors.New("unknown protocol")

	// ErrTerminated indicates the operation is terminated.
	ErrTerminated = errors.New("terminated")
)

// client implements a simple MQTT client.
type client struct {
	conn net.Conn
	r    packet.Reader
	p    Param
	log  *log.Logger

	sl   sync.Mutex // send (conn) lock
	id   uint32
	derr error

	ping  *waitop.WaitOp
	subsc *waitop.WaitOp
	unsub *waitop.WaitOp

	// message receive buffer.
	msgc *sync.Cond
	msgs []*Message
	msgr int
	msgw int

	publock sync.Mutex
}

var _ Client = (*client)(nil)

func (c *client) Disconnect(force bool) error {
	c.sl.Lock()
	defer c.sl.Unlock()
	if c.conn == nil {
		return nil
	}
	if !force {
		b, _ := (&packet.Disconnect{}).Encode()
		c.sendRaw(b)
	}
	return c.stopRaw(Explicitly)
}

func (c *client) Ping() error {
	_, err := c.ping.Do(func() error {
		return c.send(&packet.PingReq{})
	})
	return err
}

func (c *client) Subscribe(topics []Topic) error {
	var id packet.ID
	r, err := c.subsc.Do(func() error {
		array, err := packetTopics(topics)
		if err != nil {
			return err
		}
		id = c.emitID()
		return c.send(&packet.Subscribe{
			PacketID: id,
			Topics:   array,
		})
	})
	if err != nil {
		return err
	}
	p, ok := r.(*packet.SubACK)
	if !ok {
		panic(fmt.Sprintf("unexpected response: %v", r))
	}
	se := &SubscribeError{
		MismatchPacketID:    id != p.PacketID,
		MismatchResultCount: len(topics) != len(p.Results),
		RequestedQoS:        make([]QoS, len(topics)),
		ResultQoS:           make([]QoS, len(p.Results)),
	}
	for i, t := range topics {
		se.RequestedQoS[i] = t.QoS
	}
	for i, r := range p.Results {
		se.ResultQoS[i] = toQoS(r)
	}
	if se.hasErrors() {
		return se
	}
	return nil
}

func (c *client) Unsubscribe(topics []string) error {
	var id packet.ID
	r, err := c.unsub.Do(func() error {
		id = c.emitID()
		return c.send(&packet.Unsubscribe{
			PacketID: id,
			Topics:   topics,
		})
	})
	if err != nil {
		return err
	}
	p, ok := r.(*packet.UnsubACK)
	if !ok {
		panic(fmt.Sprintf("unexpected response: %v", r))
	}
	ue := &UnsubscribeError{
		MismatchPacketID: id != p.PacketID,
	}
	if ue.hasErrors() {
		return ue
	}
	return nil
}

func (c *client) Publish(qos QoS, retain bool, topic string, msg []byte) error {
	// FIXME: support AtLeastOnce and ExactlyOnce QoS
	switch qos {
	case AtMostOnce:
		return c.publish0(retain, topic, msg)
	default:
		return errors.New("unsupported QoS")
	}
}

func (c *client) Read(block bool) (*Message, error) {
	c.msgc.L.Lock()
	for c.msgr == c.msgw {
		if !block {
			c.msgc.L.Unlock()
			return nil, nil
		}
		c.msgc.Wait()
	}
	m := c.msgs[c.msgr]
	c.msgs[c.msgr] = nil
	if m != nil {
		if c.msgr++; c.msgr >= len(c.msgs) {
			c.msgr = 0
		}
	}
	c.msgc.L.Unlock()
	if m == nil {
		return nil, ErrTerminated
	}
	return m, nil
}

func (c *client) start() {
	c.ping = waitop.New()
	c.subsc = waitop.New()
	c.unsub = waitop.New()
	c.msgc = sync.NewCond(new(sync.Mutex))
	c.msgs = make([]*Message, 32)
	go c.recvLoop()
}

// stop closes connection and remove all resources.
func (c *client) stop(reason error) error {
	c.sl.Lock()
	defer c.sl.Unlock()
	return c.stopRaw(reason)
}

func (c *client) stopRaw(reason error) error {
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	c.ping.Close()
	c.subsc.Close()
	c.unsub.Close()
	if c.derr == nil {
		c.derr = reason
	}
	// clear all messages
	c.msgc.L.Lock()
	for i, m := range c.msgs {
		if m != nil {
			c.logDroppedMessage(m)
		}
		c.msgs[i] = nil
	}
	c.msgr = 0
	c.msgw = 1
	c.msgc.Signal()
	c.msgc.L.Unlock()
	return err
}

func (c *client) sendRaw(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *client) send(p packet.Packet) error {
	c.sl.Lock()
	defer c.sl.Unlock()
	if c.conn == nil {
		return errors.New("connection closed")
	}
	b, err := p.Encode()
	if err != nil {
		return err
	}
	return c.sendRaw(b)
}

func (c *client) recvLoop() {
loop:
	for {
		rp, err := packet.SplitDecode(c.r)
		if err != nil {
			nerr, ok := err.(net.Error)
			if ok && nerr.Temporary() {
				c.logTemporaryError(nerr)
				continue
			}
			c.stop(err)
			break loop
		}
		switch p := rp.(type) {
		case *packet.Publish:
			err := c.procPublish(p)
			if err != nil {
				c.stop(err)
				break loop
			}
		case *packet.SubACK:
			c.subsc.Fulfill(p)
		case *packet.UnsubACK:
			c.unsub.Fulfill(p)
		case *packet.PingResp:
			c.ping.Fulfill(p)
		default:
			c.stop(errors.New("receive unexpected packet"))
			break loop
		}
	}
	if c.p.OnDisconnect != nil {
		c.p.OnDisconnect(c.derr, c.p)
	}
	c.r = nil
}

func (c *client) publish0(retain bool, topic string, msg []byte) error {
	p := &packet.Publish{
		QoS:       AtMostOnce.qos(),
		Retain:    retain,
		TopicName: topic,
		Payload:   msg,
	}
	return c.send(p)
}

func (c *client) emitID() packet.ID {
	for {
		n := uint16(atomic.AddUint32(&c.id, 1))
		if n != 0 {
			return packet.ID(n)
		}
	}
}

func (c *client) procPublish(p *packet.Publish) error {
	// parse as Message
	var m *Message
	switch p.QoS {
	case packet.QAtMostOnce:
		m = &Message{
			Topic: p.TopicName,
			Body:  p.Payload,
		}
	default:
		// unsupported QoS.
		return errors.New("unsupported QoS")
	}
	if m == nil {
		return nil
	}
	if c.p.OnPublish != nil {
		go c.emitOnPublish(m)
		return nil
	}
	return c.put(m)
}

// put puts a message to ring buffer.
func (c *client) put(m *Message) error {
	c.msgc.L.Lock()
	c.msgs[c.msgw] = m
	if c.msgw++; c.msgw >= len(c.msgs) {
		c.msgw = 0
	}
	var dropped *Message
	if c.msgw == c.msgr {
		dropped = c.msgs[c.msgr]
		if c.msgr++; c.msgr >= len(c.msgs) {
			c.msgr = 0
		}
	}
	c.msgc.Signal()
	c.msgc.L.Unlock()
	if dropped != nil {
		c.logDroppedMessage(dropped)
	}
	return nil
}

func (c *client) emitOnPublish(m *Message) {
	c.publock.Lock()
	defer c.publock.Unlock()
	c.p.OnPublish(m)
}

func (c *client) logTemporaryError(nerr net.Error) {
	if c.log == nil {
		return
	}
	c.log.Printf("temporal error: %v", nerr)
}

func (c *client) logDroppedMessage(m *Message) {
	if c.log == nil {
		return
	}
	c.log.Printf("dropped message: %v", m)
}
