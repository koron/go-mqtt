package server

import (
	"bufio"
	"io"
	"net"
	"sync"

	"github.com/koron/go-debug"
	"github.com/koron/go-mqtt/packet"
)

// PreConn represents half-connected connection.
type PreConn interface {
	// SetData binds a user defined data to connection.
	SetData(v interface{})

	// SetReceiveHandler binds ReceiveHandler to connection.
	SetReceiveHandler(h ReceiveHandler)

	// SetSentHandler binds SentHandler to connection.
	SetSentHandler(h SentHandler)

	// SetSubscribeHandler binds SubscribeHandler to connection.
	SetSubscribeHandler(h SubscribeHandler)

	// SetPublishedHandler binds PublishedHandler to conn.
	SetPublishedHandler(h PublishedHandler)

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

	// Send sends a raw MQTT packet.
	Send(msg packet.Packet) error

	// Publish publishes a message.
	Publish(topic string, body []byte, qos packet.QoS) error

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
	subh SubscribeHandler
	pubh PublishedHandler

	wg    sync.WaitGroup
	quit  chan bool
	sendQ chan packet.Packet

	disconnect *packet.Disconnect
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
		sendQ:  make(chan packet.Packet, 1),
	}
}

func (c *conn) closeAll() {
	close(c.quit)
	close(c.sendQ)
	c.rwc.Close()
}

func (c *conn) writePacket(p packet.Packet) (int, error) {
	b, err := p.Encode()
	if err != nil {
		return 0, err
	}
	return c.writer.Write(b)
}

func (c *conn) readConnectPacket() (*packet.Connect, error) {
	b, err := packet.Split(c.reader)
	if err != nil {
		return nil, err
	}
	p := packet.Connect{}
	if err := p.Decode(b); err != nil {
		return nil, err
	}
	return &p, nil
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
	req, err := c.readConnectPacket()
	if err != nil {
		c.writePacket(&packet.ConnACK{
			ReturnCode: packet.ConnectServerUnavailable,
		})
		return err
	}
	rc := c.server.authenticate(c, req)
	// send connack packet.
	_, err = c.writePacket(&packet.ConnACK{
		SessionPresent: rc == packet.ConnectAccept,
		ReturnCode:     rc,
	})
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

func (c *conn) processMessage(p packet.Packet) error {
	switch p := p.(type) {
	case *packet.Disconnect:
		c.disconnect = p
		c.closeAll()
	case *packet.Subscribe:
		return c.processSubscribe(p)
	case *packet.Publish:
		return c.processPublish(p)
	}
	return nil
}

func (c *conn) processSubscribe(req *packet.Subscribe) error {
	resp := packet.SubACK{
		PacketID: req.PacketID,
	}
	for _, t := range req.Topics {
		// TODO: fix result determination.
		r := packet.SubscribeResult(t.RequestedQoS)
		if c.subh != nil {
			var err error
			r, err = c.subh(c, t)
			if err != nil {
				debug.Printf("mqtt: topic disabled: %s (id=%d)\n", t.Filter, c.id)
				r = packet.SubscribeFailure
			}
		}
		resp.AddResult(r)
	}
	return c.Send(&resp)
}

func (c *conn) processPublish(req *packet.Publish) error {
	if c.pubh != nil {
		err := c.pubh(c, req)
		if err != nil {
			return err
		}
	}
	// FIXME: consider PUBACK and PUBCOMP.
	return nil
}

func (c *conn) recvMain() {
loop:
	for {
		select {
		case <-c.quit:
			break loop
		default:
		}
		p, err := packet.SplitDecode(c.reader)
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
			err := c.rh(c, p)
			if err != nil {
				debug.Printf("mqtt: ReceiveHandler failed: %v (id=%d)\n", err, c.id)
				continue loop
			}
		}
		err = c.processMessage(p)
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
		case p := <-c.sendQ:
			if p == nil {
				break loop
			}
			_, err := c.writePacket(p)
			if err != nil {
				debug.Printf("mqtt: writeMessage failed: %v (id=%d)\n", err, c.id)
				continue loop
			}
			if c.sh != nil {
				err := c.sh(c, p)
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

func (c *conn) SetSubscribeHandler(h SubscribeHandler) {
	c.subh = h
}

func (c *conn) SetPublishedHandler(h PublishedHandler) {
	c.pubh = h
}

func (c *conn) Send(msg packet.Packet) error {
	// TODO: guard against sending to closed channel.
	c.sendQ <- msg
	return nil
}

func (c *conn) Publish(topic string, body []byte, qos packet.QoS) error {
	p := packet.Publish{
		Header: packet.Header{
			QoS: qos,
		},
		TopicName: topic,
		Payload:   body,
	}
	return c.Send(&p)
}

func (c *conn) Server() *Server {
	return c.server
}
