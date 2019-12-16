package client

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"

	"github.com/koron/go-mqtt/internal/waitop"
	"github.com/koron/go-mqtt/packet"
	"golang.org/x/net/websocket"
)

// Connect connects to MQTT broker and returns a Client.
func Connect(p Param) (Client, error) {
	c, err := dial(p)
	if err != nil {
		return nil, err
	}
	r := p.newPacketReader(c)

	// send CONNECT packet.
	bc, err := p.connectPacket().Encode()
	if err != nil {
		c.Close()
		return nil, err
	}
	_, err = c.Write(bc)
	if err != nil {
		c.Close()
		return nil, err
	}

	// receive CONNACK packet.
	rp, err := packet.SplitDecode(r)
	if err != nil {
		c.Close()
		return nil, err
	}
	ack, ok := rp.(*packet.ConnACK)
	if !ok {
		c.Close()
		return nil, errors.New("received non CONNACK")
	}
	if ack.ReturnCode != packet.ConnectAccept {
		return nil, ack.ReturnCode
	}

	opts := p.options()
	cl := &client{
		conn: c,
		quit: make(chan bool, 1),
		r:    r,
		p:    p,
		log:  opts.Logger,
		kd:   opts.keepAliveInterval(),
		wt:   map[packet.ID]*waitop.WaitOp{},
	}
	cl.start()
	return cl, nil
}

func dial(p Param) (net.Conn, error) {
	u, err := p.url()
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "tcp":
		return dialTCP(p, u)
	case "ssl", "tcps", "tls":
		return dialTLS(p, u)
	case "ws":
		return dialWS(p, u)
	case "wss":
		return dialWSS(p, u)
	}
	return nil, ErrUnknownProtocol
}

func dialTCP(p Param, u *url.URL) (net.Conn, error) {
	opts := p.options()
	c, err := net.DialTimeout("tcp", u.Host, opts.ConnectTimeout)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func dialTLS(p Param, u *url.URL) (net.Conn, error) {
	opts := p.options()
	c, err := tls.DialWithDialer(opts.dialer(), "tcp", u.Host, opts.TLSConfig)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func dialWS(p Param, u *url.URL) (net.Conn, error) {
	opts := p.options()
	cnf, err := websocket.NewConfig(u.String(), opts.wsOrigin(u))
	if err != nil {
		return nil, err
	}
	cnf.Dialer = opts.dialer()
	c, err := websocket.DialConfig(cnf)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func dialWSS(p Param, u *url.URL) (net.Conn, error) {
	opts := p.options()
	cnf, err := websocket.NewConfig(u.String(), opts.wsOrigin(u))
	if err != nil {
		return nil, err
	}
	cnf.Dialer = opts.dialer()
	cnf.TlsConfig = opts.TLSConfig
	c, err := websocket.DialConfig(cnf)
	if err != nil {
		return nil, err
	}
	return c, nil
}
