package server

import (
	"bufio"
	"io"
	"net"

	"github.com/koron/go-debug"
	"github.com/surgemq/message"
)

type connID uint64

type conn struct {
	server *Server
	rwc    net.Conn
	reader *bufio.Reader
	writer io.Writer
	id     connID
}

func newConn(srv *Server, rwc net.Conn) *conn {
	return &conn{
		server: srv,
		rwc:    rwc,
		reader: bufio.NewReader(rwc),
		writer: rwc,
	}
}

func (c *conn) Close() error {
	c.rwc.Close()
	// TODO: terminate goroutines for a conn.
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
	defer c.Close()
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
	go c.recvMain()
	c.sendMain()
	c.server.unregister(c)
}

func (c *conn) recvMain() {
	for {
		// TODO:
	}
}

func (c *conn) sendMain() {
	for {
		// TODO:
	}
}
