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

	ConnectTimeout time.Duration
	TLSConfig *tls.Config

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

// DefaultOptions represents default values which used for when Connect()'s
// opts argument is nil.
var DefaultOptions = &Options{
	Version:      4,
	CleanSession: true,
	KeepAlive:    60,
}
