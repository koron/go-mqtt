package client

import (
	"errors"
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

	// ReadMessage returns a message if it was available.  Otherwise this will
	// block.
	ReadMessage() (*Message, error)

	// PeekMessage returns true if ReadMessage() can return one or more
	// messages without blocking.
	PeekMessage() bool
}

var (
	// ErrTerminated indicates the operation is terminated.
	ErrTerminated = errors.New("ping terminated")
)

// client implements a simple MQTT client.
type client struct {
	conn net.Conn
	r    packet.Reader
	sl   sync.Mutex // send (conn) lock

	ping *waitop.WaitOp

	id uint32
}

var _ Client = (*client)(nil)

func (c *client) Disconnect(force bool) error {
	if c.conn == nil {
		return nil
	}
	if !force {
		b, _ := (&packet.Disconnect{}).Encode()
		c.send(b)
	}
	// close connection and remove all resources.
	c.sl.Lock()
	err := c.conn.Close()
	c.conn = nil
	c.r = nil
	c.ping.Close()
	c.sl.Unlock()
	return err
}

func (c *client) Ping() error {
	_, err := c.ping.Do(func() error {
		return c.encodeAndSend(&packet.PingReq{})
	})
	return err
}

func (c *client) Subscribe(topics []Topic) error {
	// TODO:
	return nil
}

func (c *client) Unsubscribe(topics []string) error {
	// TODO:
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

func (c *client) PeekMessage() bool {
	// TODO:
	return false
}

func (c *client) ReadMessage() (*Message, error) {
	// TODO:
	return nil, nil
}

func (c *client) start() {
	c.ping = waitop.New()
	// FIXME: initialize client's other fields
	go c.recvLoop()
}

func (c *client) stop(err error) {
	// TODO: disconnect and stop
}

func (c *client) logTemporaryError(err error) {
	// TODO:
}

func (c *client) send(b []byte) error {
	c.sl.Lock()
	defer c.sl.Unlock()
	if c.conn == nil {
		return errors.New("connection closed")
	}
	_, err := c.conn.Write(b)
	return err
}

func (c *client) encodeAndSend(p packet.Packet) error {
	b, err := p.Encode()
	if err != nil {
		return err
	}
	return c.send(b)
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
			// TODO:
		case *packet.SubACK:
			// TODO:
		case *packet.UnsubACK:
			// TODO:
		case *packet.PingResp:
			c.ping.Fulfill(p)
		default:
			c.stop(errors.New("receive unexpected packet"))
			break loop
		}
	}
}

func (c *client) publish0(retain bool, topic string, msg []byte) error {
	p := &packet.Publish{
		QoS:       AtMostOnce.qos(),
		Retain:    retain,
		TopicName: topic,
		Payload:   msg,
	}
	err := c.encodeAndSend(p)
	if err != nil {
		// TODO: treat temporary error.
		return err
	}
	return nil
}

func (c *client) emitID() packet.ID {
	n := atomic.AddUint32(&c.id, 1)
	return packet.ID(n)
}
