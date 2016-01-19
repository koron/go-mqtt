package server

import (
	"bufio"
	"net"
	"time"

	"github.com/koron/go-mqtt/internal/backoff"
	"github.com/koron/go-mqtt/packet"
)

type client struct {
	srv  *Server
	conn net.Conn
	rd   packet.Reader
	ca   ClientAdapter
}

func newClient(srv *Server, conn net.Conn) *client {
	return &client{
		srv:  srv,
		conn: conn,
		rd:   bufio.NewReader(conn),
	}
}

func (c *client) serve() {
	end := func() {
		c.conn.Close()
		c.srv.clientOnDisconnect(c)
	}
	err := c.establish()
	if err != nil {
		c.srv.logEstablishFailure(c, err)
		end()
		return
	}
	c.srv.clientOnStart(c)
	err = c.recvLoop()
	c.srv.clientOnStop(c, err)
	end()
}

func (c *client) establish() error {
	p, err := c.receiveConnect()
	if err != nil {
		return err
	}
	c.ca, err = c.srv.clientOnConnect(c, p)
	if err != nil {
		rc := packet.ConnectNotAuthorized
		if cerr, ok := err.(ConnectError); ok {
			rc = cerr.toRC()
		}
		c.send(&packet.ConnACK{ReturnCode: rc})
		return err
	}
	// send success ConnACK.
	err = c.send(&packet.ConnACK{
		SessionPresent: c.ca.IsSessionPresent(),
		ReturnCode:     packet.ConnectAccept,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) receiveConnect() (*packet.Connect, error) {
	b, err := packet.Split(c.rd)
	if err != nil {
		return nil, err
	}
	p := new(packet.Connect)
	if err := p.Decode(b); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *client) recvLoop() error {
	delay := backoff.Exp{Min: time.Millisecond * 5}
	for {
		p, err := packet.SplitDecode(c.rd)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				c.srv.logTemporaryError(nerr, &delay, c)
				delay.Wait()
				continue
			}
			return err
		}
		delay.Reset()
		err = c.process(p)
		if err != nil {
			return err
		}
	}
}

func (c *client) process(raw packet.Packet) error {
	switch p := raw.(type) {
	// TODO:
	}
	return nil
}

func (c *client) send(p packet.Packet) error {
	b, err := p.Encode()
	if err != nil {
		return err
	}
	_, err = c.conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}
