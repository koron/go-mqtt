package client

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/koron/go-mqtt/packet"
)

// Param represents connection parameters for MQTT client.
type Param struct {
	// Addr is URL to connect like "tcp://192.168.0.1:1883".
	Addr string

	// ID is used as MQTT's client ID.
	ID string

	// OnPublish is called when receive a PUBLISH message with independed
	// goroutine.  If it is omitted, received messages are stored into buffer.
	OnPublish PublishedFunc

	// OnDisconnect is called when connection is disconnected.
	OnDisconnect DisconnectedFunc

	// Options is option parameters for connection.
	Options *Options
}

func (p *Param) options() *Options {
	if p.Options == nil {
		return DefaultOptions
	}
	return p.Options
}

func (p *Param) addr() string {
	if p.Addr == "" {
		return "tcp://127.0.0.1:1883"
	}
	return p.Addr
}

func (p *Param) url() (*url.URL, error) {
	return url.Parse(p.addr())
}

func (p *Param) connectPacket() *packet.Connect {
	return p.options().connectPacket(p.ID)
}

func (p *Param) newPacketReader(c net.Conn) packet.Reader {
	// TODO: apply timeout configuration.
	r := bufio.NewReader(c)
	return r
}

// Options represents connect options
type Options struct {
	Version      uint8   // MQTT's protocol version 3 or 4 (fallback to 4)
	Username     *string // username to connect (option)
	Password     *string // password to connect (option)
	CleanSession bool
	KeepAlive    uint16
	Will         *Will

	// DisableAutoKeepAlive disables auto ping to keep alive.
	DisableAutoKeepAlive bool

	ConnectTimeout time.Duration
	TLSConfig      *tls.Config

	WSOrigin string

	Logger *log.Logger
}

func (o *Options) version() uint8 {
	switch o.Version {
	case 3:
		return 3
	default:
		return 4
	}
}

func (o *Options) connectPacket(id string) *packet.Connect {
	p := &packet.Connect{
		ClientID:     id,
		Version:      o.version(),
		Username:     o.Username,
		Password:     o.Password,
		CleanSession: o.CleanSession,
		KeepAlive:    o.KeepAlive,
	}
	if o.Will != nil {
		p.WillFlag = true
		p.WillQoS = o.Will.QoS.qos()
		p.WillRetain = o.Will.Retain
		p.WillTopic = o.Will.Topic
		p.WillMessage = o.Will.Message
	}
	return p
}

func (o *Options) keepAliveInterval() time.Duration {
	const faster = time.Millisecond * 500
	d := time.Second * time.Duration(o.KeepAlive)
	if d <= faster {
		return d
	}
	return d - faster
}

func (o *Options) dialer() *net.Dialer {
	return &net.Dialer{Timeout: o.ConnectTimeout}
}

func (o *Options) wsOrigin(u *url.URL) string {
	if o.WSOrigin != "" {
		return o.WSOrigin
	}
	return wsOrigin(u)
}

func wsOrigin(u *url.URL) string {
	v := *u
	// convert schema.
	if v.Scheme == "wss" {
		v.Scheme = "https"
	} else {
		v.Scheme = "http"
	}
	// keep user and host parts, then reset other parts.
	v.Opaque = ""
	v.Path = ""
	v.ForceQuery = false
	v.RawQuery = ""
	v.Fragment = ""
	return v.String()
}

// DefaultOptions represents default values which used for when Connect()'s
// opts argument is nil.
var DefaultOptions = &Options{
	Version:      4,
	CleanSession: true,
	KeepAlive:    30,
}
