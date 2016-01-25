package server

import (
	"bufio"
	"net"
	"sync"
	"time"

	"github.com/koron/go-mqtt/internal/backoff"
	"github.com/koron/go-mqtt/packet"
)

type client struct {
	srv  *Server
	conn net.Conn
	wg   sync.WaitGroup
	quit chan bool
	sq   chan packet.Packet
	rd   packet.Reader
	ca   ClientAdapter
	pf   PacketFilter
}

func newClient(srv *Server, conn net.Conn) *client {
	return &client{
		srv:  srv,
		conn: conn,
		quit: make(chan bool, 1),
		sq:   make(chan packet.Packet, 1),
		rd:   bufio.NewReader(conn),
	}
}

func (c *client) terminate() {
	close(c.quit)
	c.conn.Close()
}

func (c *client) id() string {
	return c.ca.ID()
}

func (c *client) serve() {
	err := c.establish()
	if err != nil {
		close(c.quit)
		c.conn.Close()
		close(c.sq)
		c.srv.clientOnDisconnect(c, err)
		return
	}
	c.srv.clientOnStart(c)
	c.wg.Add(1)
	go c.sendLoop()
	err = c.recvLoop()
	close(c.sq)
	c.wg.Wait() // wait to terminate sendLoop
	c.srv.clientOnStop(c)
	c.srv.clientOnDisconnect(c, err)
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
	if pf, ok := c.ca.(PacketFilter); ok {
		c.pf = pf
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

func (c *client) sendLoop() {
	defer c.wg.Done()
	for {
		select {
		case <-c.quit:
			return
		case p := <-c.sq:
			if p == nil {
				return
			}
			err := c.send(p)
			if err != nil {
				c.srv.logSendPacketError(c, p, err)
			}
		}
	}
}

func (c *client) recvLoop() error {
	delay := backoff.Exp{Min: time.Millisecond * 5}
	for {
		p, err := packet.SplitDecode(c.rd)
		select {
		case <-c.quit:
			return nil
		default:
		}
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
			if cerr, ok := err.(*clientError); ok && cerr.cont {
				if cerr == errDisconnected {
					return nil
				}
				c.srv.logClientError(cerr, p, c)
				continue
			}
			return err
		}
	}
}

func (c *client) process(raw packet.Packet) error {
	if c.pf != nil {
		err := c.pf.PreProcess(raw)
		if err != nil {
			return err
		}
	}
	switch p := raw.(type) {
	case *packet.Disconnect:
		return c.processDisconnect(p)
	case *packet.PingReq:
		return c.processPingReq(p)
	case *packet.Subscribe:
		return c.processSubscribe(p)
	case *packet.Unsubscribe:
		return c.processUnsubscribe(p)
	case *packet.Publish:
		return c.processPublish(p)
	case *packet.PubACK:
		return c.processPubACK(p)
	case *packet.PubRec:
		return c.processPubRec(p)
	case *packet.PubRel:
		return c.processPubRel(p)
	case *packet.PubComp:
		return c.processPubComp(p)
	default:
		return errNotAcceptable
	}
}

func (c *client) processDisconnect(p *packet.Disconnect) error {
	err := c.ca.OnDisconnect()
	if err != nil {
		return err
	}
	return errDisconnected
}

func (c *client) processPingReq(p *packet.PingReq) error {
	f, err := c.ca.OnPing()
	if err != nil {
		return err
	}
	if f {
		c.sq <- &packet.PingResp{}
	}
	return nil
}

func (c *client) processSubscribe(p *packet.Subscribe) error {
	// TODO:
	return errNotSuported
}

func (c *client) processUnsubscribe(p *packet.Unsubscribe) error {
	// TODO:
	return errNotSuported
}

func (c *client) processPublish(p *packet.Publish) error {
	// TODO:
	return errNotSuported
}

func (c *client) processPubACK(p *packet.PubACK) error {
	// FIXME: QoS1 will be supported in future.
	return errNotSuported
}

func (c *client) processPubRec(p *packet.PubRec) error {
	// FIXME: QoS2 will be supported in future.
	return errNotSuported
}

func (c *client) processPubRel(p *packet.PubRel) error {
	// FIXME: QoS2 will be supported in future.
	return errNotSuported
}

func (c *client) processPubComp(p *packet.PubComp) error {
	// FIXME: QoS2 will be supported in future.
	return errNotSuported
}

func (c *client) send(p packet.Packet) error {
	b, err := p.Encode()
	if err != nil {
		return err
	}
	if c.pf == nil {
		// send without PacketFilter
		_, err = c.conn.Write(b)
		if err != nil {
			return err
		}
		return nil
	}
	// send with PacketFilter
	b2, err := c.pf.PreSend(p, b)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(b2)
	if err != nil {
		return err
	}
	c.pf.PostSend(p, b2)
	return nil
}
