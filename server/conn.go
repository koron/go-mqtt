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
	// SetReceiveHandler binds ReceiveHandler to connection.
	SetReceiveHandler(rh ReceiveHandler)

	// Conn returns corresponding net.Conn.
	Conn() net.Conn
}

// Conn represents a MQTT client connection.
type Conn interface {
	// Close closes a connection.
	Close() error

	// Send sends a message to connection.
	Send(msg message.Message) error

	// Server returns corresponding server.
	Server() *Server

	// Conn returns corresponding net.Conn.
	Conn() net.Conn
}

type connID uint64

type conn struct {
	server *Server
	rwc    net.Conn
	reader *bufio.Reader
	writer io.Writer

	id    connID
	wg    sync.WaitGroup
	quit  chan bool
	sendQ chan message.Message
	rh    ReceiveHandler
}

var (
	_ Conn    = (*conn)(nil)
	_ PreConn = (*conn)(nil)
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

func (c *conn) Close() error {
	close(c.quit)
	close(c.sendQ)
	c.rwc.Close()
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
			debug.Printf("mqtt: readMessage failed: %v, id=%d\n", err, c.id)
			continue
		}
		if c.rh != nil {
			err := c.rh(c, msg)
			if err != nil {
				debug.Printf("mqtt: ReceiveHandler failed: %v, id=%d\n", err, c.id)
			}
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
			_, err := writeMessage(c.writer, m)
			if err != nil {
				debug.Printf("mqtt: writeMessage failed: %v, id=%d\n", err, c.id)
			}
		}
	}
	c.wg.Done()
}

func (c *conn) Conn() net.Conn {
	return c.rwc
}

func (c *conn) SetReceiveHandler(rh ReceiveHandler) {
	c.rh = rh
}

func (c *conn) Send(msg message.Message) error {
	// FIXME: guard against sending to closed channel.
	c.sendQ <- msg
	return nil
}

func (c *conn) Server() *Server {
	return c.server
}
