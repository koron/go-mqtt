package client

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/koron/go-mqtt/packet"
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

	cl := &client{
		conn: c,
		r:    r,
		p:    p,
		log:  p.options().Logger,
	}
	cl.start()
	return cl, nil
}

func dial(p Param) (net.Conn, error) {
	u, err := p.url()
	if err != nil {
		return nil, err
	}
	opts := p.options()
	to := opts.ConnectTimeout
	switch u.Scheme {
	case "tcp":
		c, err := net.DialTimeout("tcp", u.Host, to)
		if err != nil {
			return nil, err
		}
		return c, nil
	case "ssl", "tcps", "tls":
		c, err := tls.DialWithDialer(&net.Dialer{Timeout: to},
			"tcp", u.Host, opts.TLSConfig)
		if err != nil {
			return nil, err
		}
		return c, nil
	default:
		return nil, ErrUnknownProtocol
	}
}
