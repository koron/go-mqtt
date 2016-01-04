package server

import (
	"bufio"
	"io"
	"net"
	"sync"

	"github.com/koron/go-debug"
	"github.com/surgemq/message"
)

// Conn represents a MQTT client connection.
type Conn interface {
}

type connID uint64

type conn struct {
	id connID

	server *Server
	rwc    net.Conn
	reader *bufio.Reader
	writer io.Writer

	quit chan bool
	wg   sync.WaitGroup
}

func newConn(srv *Server, rwc net.Conn) *conn {
	return &conn{
		server: srv,
		rwc:    rwc,
		reader: bufio.NewReader(rwc),
		writer: rwc,
		quit:   make(chan bool, 1),
	}
}

func (c *conn) Close() error {
	close(c.quit)
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
	err = writeMessage(c.writer, resp)
	if err != nil {
		return err
	}
	return nil
}

func (c *conn) serve() {
	defer c.rwc.Close()
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
		}
		// TODO:
	}
	c.wg.Done()
}

func (c *conn) sendMain() {
loop:
	for {
		select {
		case <-c.quit:
			break loop
		}
		// TODO:
	}
	c.wg.Done()
}
