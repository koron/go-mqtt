package server

import (
	"bufio"
	"io"
	"net"
	"sync"

	"github.com/koron/go-debug"
	"github.com/surgemq/message"
)

// PreConn represents half-connected connection.
type PreConn interface {
	// SetData binds a user defined data to connection.
	SetData(v interface{})

	// SetReceiveHandler binds ReceiveHandler to connection.
	SetReceiveHandler(h ReceiveHandler)

	// SetSentHandler binds SentHandler to connection.
	SetSentHandler(h SentHandler)

	// SetSubscribeHandler binds SubscribleHandler to connection.
	SetSubscribeHandler(h SubscribleHandler)

	// Conn returns corresponding net.Conn.
	Conn() net.Conn
}

// Conn represents a MQTT client connection.
type Conn interface {
	// ID returns connection ID.
	ID() ConnID

	// Data returns user corresponded data.
	Data() interface{}

	// Server returns corresponding server.
	Server() *Server

	// Conn returns corresponding net.Conn.
	Conn() net.Conn

	// Send sends a raw MQTT message.
	Send(msg message.Message) error

	// Publish publishes a message.
	Publish(topic string, body []byte, qos byte) error

	// Close closes a connection.
	Close() error
}

// DisConn represents a closed MQTT client connection.
type DisConn interface {
	// ID returns connection ID.
	ID() ConnID

	// Data returns user corresponded data.
	Data() interface{}
}

// ConnID identifies connections in a Server.
type ConnID uint64

type conn struct {
	server *Server
	rwc    net.Conn
	reader *bufio.Reader
	writer io.Writer

	id   ConnID
	data interface{}
	rh   ReceiveHandler
	sh   SentHandler
	subh SubscribleHandler

	wg    sync.WaitGroup
	quit  chan bool
	sendQ chan message.Message

	disconnect *message.DisconnectMessage
}

var (
	_ Conn    = (*conn)(nil)
	_ PreConn = (*conn)(nil)
	_ DisConn = (*conn)(nil)
)

func newConn(srv *Server, rwc net.Conn) *conn {
	return &conn{
		server: srv,
		rwc:    rwc,
		reader: bufio.NewReader(rwc),
		writer: rwc,
		quit:   make(chan bool, 1),
		sendQ:  make(chan message.Message, 1),
	}
}

func (c *conn) closeAll() {
	close(c.quit)
	close(c.sendQ)
	c.rwc.Close()
}

func (c *conn) ID() ConnID {
	return c.id
}

func (c *conn) Data() interface{} {
	return c.data
}

func (c *conn) SetData(v interface{}) {
	c.data = v
}

func (c *conn) Close() error {
	c.closeAll()
	c.wg.Wait()
	return nil
}

func (c *conn) establishConnection() error {
	req, err := readConnectMessage(c.reader)
	if err != nil {
		writeConnackErrorMessage(c.writer, err)
		return err
	}
	err = c.server.authenticate(c, req)
	if err != nil {
		writeConnackErrorMessage(c.writer, err)
		return err
	}
	// send connack message.
	resp := message.NewConnackMessage()
	resp.SetSessionPresent(true)
	resp.SetReturnCode(message.ConnectionAccepted)
	_, err = writeMessage(c.writer, resp)
	if err != nil {
		return err
	}
	return nil
}

func (c *conn) serve() {
	err := c.establishConnection()
	if err != nil {
		debug.Printf("mqtt: establishConnection failed: %v\n", err)
		return
	}
	err = c.server.register(c)
	if err != nil {
		debug.Printf("mqtt: register failed: %v\n", err)
		return
	}
	c.wg = sync.WaitGroup{}
	c.wg.Add(2)
	go c.recvMain()
	c.sendMain()
	c.server.unregister(c)
}

func (c *conn) processMessage(msg message.Message) error {
	switch m := msg.(type) {
	case *message.DisconnectMessage:
		c.disconnect = m
		c.closeAll()
	case *message.SubscribeMessage:
		return c.processSubscribe(m)
	}
	return nil
}

func (c *conn) processSubscribe(req *message.SubscribeMessage) error {
	topics := req.Topics()
	qos := req.Qos()
	rc := make([]byte, 0, len(topics))
	for i, t := range topics {
		tqos := qos[i]
		if c.subh != nil {
			q, err := c.subh(c, string(t), qos[i])
			if err != nil {
				debug.Printf("mqtt: topic disabled: %s (id=%d)\n", t, c.id)
				q = message.QosFailure
			}
			tqos = q
		}
		rc = append(rc, tqos)
	}
	// send SUBACK
	resp := message.NewSubackMessage()
	resp.SetPacketId(req.PacketId())
	if err := resp.AddReturnCodes(rc); err != nil {
		return err
	}
	return c.Send(resp)
}

func (c *conn) recvMain() {
loop:
	for {
		select {
		case <-c.quit:
			break loop
		default:
		}
		msg, err := readMessage(c.reader)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				// TODO: exponential backoff sleep.
				continue
			}
			debug.Printf("mqtt: conn closed: %v (id=%d)", err, c.id)
			c.closeAll()
			break loop
		}
		if c.rh != nil {
			err := c.rh(c, msg)
			if err != nil {
				debug.Printf("mqtt: ReceiveHandler failed: %v (id=%d)\n", err, c.id)
				continue loop
			}
		}
		err = c.processMessage(msg)
		if err != nil {
			debug.Printf("mqtt: processMessage failed: %v (id=%d)\n", err, c.id)
		}
	}
	c.wg.Done()
}

func (c *conn) sendMain() {
loop:
	for {
		select {
		case <-c.quit:
			break loop
		case m := <-c.sendQ:
			if m == nil {
				break loop
			}
			_, err := writeMessage(c.writer, m)
			if err != nil {
				debug.Printf("mqtt: writeMessage failed: %v (id=%d)\n", err, c.id)
				continue loop
			}
			if c.sh != nil {
				err := c.sh(c, m)
				if err != nil {
					debug.Printf("mqtt: SentHandler failed %v (id=%d)\n", err, c.id)
				}
			}
		}
	}
	c.wg.Done()
}

func (c *conn) Conn() net.Conn {
	return c.rwc
}

func (c *conn) SetReceiveHandler(h ReceiveHandler) {
	c.rh = h
}

func (c *conn) SetSentHandler(h SentHandler) {
	c.sh = h
}

func (c *conn) SetSubscribeHandler(h SubscribleHandler) {
	c.subh = h
}

func (c *conn) Send(msg message.Message) error {
	// TODO: guard against sending to closed channel.
	c.sendQ <- msg
	return nil
}

func (c *conn) Publish(topic string, body []byte, qos byte) error {
	msg := message.NewPublishMessage()
	msg.SetTopic([]byte(topic))
	msg.SetQoS(qos)
	msg.SetPayload(body)
	return c.Send(msg)
}

func (c *conn) Server() *Server {
	return c.server
}
