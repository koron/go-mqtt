package server

import "github.com/koron/go-mqtt/packet"

// Message represents published message.
type Message struct {
	Dup    bool
	QoS    QoS
	Retain bool
	Topic  string
	Body   []byte
}

func toMessage(p *packet.Publish) *Message {
	return &Message{
		Dup:    p.Dup,
		QoS:    toQoS(p.QoS),
		Retain: p.Retain,
		Topic:  p.TopicName,
		Body:   p.Payload,
	}
}
