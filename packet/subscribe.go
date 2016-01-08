package packet

import (
	"bytes"
	"fmt"
)

// Subscribe represents SUBSRIBE packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#subscribe
type Subscribe struct {
	Header
	MessageID MessageID
	Topics    []Topic
}

var _ Packet = (*Subscribe)(nil)

// Encode returns serialized Subscribe packet.
func (p *Subscribe) Encode() ([]byte, error) {
	var (
		header = &Header{
			Type: TSubscribe,
			Dup:  p.Dup,
			QoS:  p.QoS,
		}
		messageID = p.MessageID.bytes()
		topics    []byte
	)
	topics, err := encodeTopics(p.Topics)
	if err != nil {
		return nil, err
	}
	return encode(header, messageID, topics)
}

// Decode deserializes []byte as Subscribe packet.
func (p *Subscribe) Decode(b []byte) error {
	// TODO: implement me.
	return nil
}

// AddTopic adds a topic to SUBSRIBE packet.
func (p *Subscribe) AddTopic(topic Topic) {
	if p.Topics == nil {
		p.Topics = make([]Topic, 0, 4)
	}
	p.Topics = append(p.Topics, topic)
}

// SubACK represents SUBACK packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#suback
type SubACK struct {
	Header
	MessageID        MessageID
	GrantedQoSLevels []QoS
}

var _ Packet = (*SubACK)(nil)

// Encode returns serialized SubACK packet.
func (p *SubACK) Encode() ([]byte, error) {
	// a vector of granted QoS levels.
	b := make([]byte, len(p.GrantedQoSLevels))
	for i, qos := range p.GrantedQoSLevels {
		b[i] = byte(qos & 0x3)
	}
	return encode(&Header{Type: TSubACK}, p.MessageID.bytes(), b)
}

// Decode deserializes []byte as SubACK packet.
func (p *SubACK) Decode(b []byte) error {
	// TODO: implement me.
	return nil
}

// Unsubscribe represents UNSUBSCRIBE packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#unsubscribe
type Unsubscribe struct {
	Header
	MessageID MessageID
	Topics    []string
}

var _ Packet = (*Unsubscribe)(nil)

// Encode returns serialized Unsubscribe packet.
func (p *Unsubscribe) Encode() ([]byte, error) {
	var (
		header = &Header{
			Type: TUnsubscribe,
			Dup:  p.Dup,
			QoS:  p.QoS,
		}
		messageID = p.MessageID.bytes()
		topics    bytes.Buffer
	)
	for i, t := range p.Topics {
		b := encodeString(t)
		if b == nil {
			return nil, fmt.Errorf("too long topic name in #%d", i)
		}
		_, err := topics.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return encode(header, messageID, topics.Bytes())
}

// Decode deserializes []byte as Unsubscribe packet.
func (p *Unsubscribe) Decode(b []byte) error {
	// TODO: implement me.
	return nil
}

// UnsubACK represents UNSUBACK packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#unsuback
type UnsubACK struct {
	Header
	MessageID MessageID
}

var _ Packet = (*UnsubACK)(nil)

// Encode returns serialized UnsubACK packet.
func (p *UnsubACK) Encode() ([]byte, error) {
	return encode(&Header{Type: TUnsubACK}, p.MessageID.bytes())
}

// Decode deserializes []byte as UnsubACK packet.
func (p *UnsubACK) Decode(b []byte) error {
	// TODO: implement me.
	return nil
}

// Topic represents topics to subscribe.
type Topic struct {
	Name string
	QoS  QoS
}

func encodeTopics(topics []Topic) ([]byte, error) {
	buf := bytes.Buffer{}
	for i, t := range topics {
		n := encodeString(t.Name)
		if n == nil {
			return nil, fmt.Errorf("too long topic name in #%d", i)
		}
		_, err := buf.Write(n)
		if err != nil {
			return nil, err
		}
		err = buf.WriteByte(byte(t.QoS & 0x03))
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
