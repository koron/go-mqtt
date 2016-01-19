package server

import "github.com/koron/go-mqtt/packet"

// ReceiveHandler called when receive a MQTT packet.
type ReceiveHandler func(conn Conn, p packet.Packet) error

// SentHandler called after sent a MQTT packet.
type SentHandler func(conn Conn, p packet.Packet) error

// SubscribeHandler called each topic to subscribe.
type SubscribeHandler func(conn Conn, t packet.Topic) (r packet.SubscribeResult, err error)

// PublishedHandler notifies a PUBLISH packet.
type PublishedHandler func(con Conn, p *packet.Publish) error
