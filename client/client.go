package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koron/go-mqtt/internal/backoff"
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
	wg   sync.WaitGroup
	quit chan bool
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

	// auto keep aliving
	kd time.Duration
	kl sync.Mutex
	kx chan struct{}

	wl sync.RWMutex
	wt map[packet.ID]*waitop.WaitOp
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
	// FIXME: support ExactlyOnce QoS
	switch qos {
	case AtMostOnce:
		return c.publish0(retain, topic, msg)
	case AtLeastOnce:
		return c.Publish1(context.Background(), retain, topic, msg)
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
	if !c.p.options().DisableAutoKeepAlive {
		go c.keepAliveLoop()
	}
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
	close(c.quit)
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
	err = c.sendRaw(b)
	if err != nil {
		return err
	}
	c.keepAliveExtend()
	return nil
}

func (c *client) keepAliveLoop() {
	c.kl.Lock()
	c.kx = make(chan struct{})
	ti := time.NewTimer(c.kd)
	tistop := func() {
		if !ti.Stop() {
			<-ti.C
		}
	}
	c.kl.Unlock()
loop:
	for {
		select {
		case <-c.quit:
			break loop
		case <-c.kx:
			tistop()
			ti.Reset(c.kd)
		case <-ti.C:
			go c.Ping()
			// c.keepAliveExtend() will be called and resumed the Timer by
			// c.send() when sending Ping packet.
		}
	}
	c.kl.Lock()
	tistop()
	close(c.kx)
	c.kx = nil
	c.kl.Unlock()
}

func (c *client) keepAliveExtend() {
	c.kl.Lock()
	if c.kx != nil {
		select {
		case c.kx <- struct{}{}:
		default:
		}
	}
	c.kl.Unlock()
}

func (c *client) recvLoop() {
	delay := backoff.Exp{Min: time.Millisecond * 5}
loop:
	for {
		p, err := packet.SplitDecode(c.r)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				c.logTemporaryError(nerr)
				delay.Wait()
				continue
			}
			c.stop(err)
			break loop
		}
		delay.Reset()
		if err := c.dispatch(p); err != nil {
			c.stop(err)
			break loop
		}
	}
	if c.p.OnDisconnect != nil {
		c.p.OnDisconnect(c.derr, c.p)
	}
	c.r = nil
}

// dispatch dispatches received packet.
func (c *client) dispatch(raw packet.Packet) error {
	switch p := raw.(type) {
	case *packet.Publish:
		return c.procPublish(p)
	case *packet.PubACK:
		c.doneWaitOp(p.PacketID)
	case *packet.SubACK:
		c.subsc.Fulfill(p)
	case *packet.UnsubACK:
		c.unsub.Fulfill(p)
	case *packet.PingResp:
		c.ping.Fulfill(p)
	default:
		return errors.New("receive unexpected packet")
	}
	return nil
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

// Publish1 publishes a message with QoS=1 (at least once). This blocks until
// receive PubACK or context is exceeded.
func (c *client) Publish1(ctx context.Context, retain bool, topic string, msg []byte) error {
	id := c.emitID()
	w, err := c.newWaitOp(id)
	if err != nil {
		return err
	}
	defer c.closeWaitOp(id)
	// FIXME: support context
	_, err = w.Do(func() error {
		return c.send(&packet.Publish{
			QoS:       AtLeastOnce.qos(),
			Retain:    retain,
			TopicName: topic,
			PacketID:  id,
			Payload:   msg,
		})
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) newWaitOp(id packet.ID) (*waitop.WaitOp, error) {
	c.wl.Lock()
	defer c.wl.Unlock()
	if _, ok := c.wt[id]; ok {
		return nil, fmt.Errorf("duplicated wait ID: %d", id)
	}
	w := waitop.New()
	c.wt[id] = w
	return w, nil
}

func (c *client) doneWaitOp(id packet.ID) {
	c.wl.RLock()
	defer c.wl.RUnlock()
	w, ok := c.wt[id]
	if !ok {
		// FIXME: log ignore Packet ID.
		return
	}
	w.Fulfill(id)
}

func (c *client) closeWaitOp(id packet.ID) {
	c.wl.Lock()
	delete(c.wt, id)
	c.wl.Unlock()
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
