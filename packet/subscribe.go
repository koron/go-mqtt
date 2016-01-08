package packet

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
	// TODO: rewrite with newDecoder()
	if len(b) < 2 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TSubscribe {
		return errors.New("type mismatch")
	}
	data, err := decodeRemain(b[2:])
	if err != nil {
		return nil
	}
	if len(data) < 2 {
		return errors.New("too short remain")
	}
	topics, err := decodeTopics(bytes.NewReader(data[2:]))
	if err != nil {
		return err
	}
	p.Header.decode(b[0])
	p.MessageID = decodeMessageID(data[0:2])
	p.Topics = topics
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
	MessageID MessageID
	Results   []SubscribeResult
}

var _ Packet = (*SubACK)(nil)

// Encode returns serialized SubACK packet.
func (p *SubACK) Encode() ([]byte, error) {
	// a vector of granted QoS levels.
	b := make([]byte, len(p.Results))
	for i, r := range p.Results {
		b[i] = byte(r)
	}
	return encode(&Header{Type: TSubACK}, p.MessageID.bytes(), b)
}

// Decode deserializes []byte as SubACK packet.
func (p *SubACK) Decode(b []byte) error {
	// TODO: rewrite with newDecoder()
	if len(b) < 2 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TSubACK {
		return errors.New("type mismatch")
	}
	data, err := decodeRemain(b[2:])
	if err != nil {
		return nil
	}
	if len(data) < 2 {
		return errors.New("too short remain")
	}
	results := make([]SubscribeResult, len(data)-2)
	for i, d := range data[2:] {
		switch d {
		case 0x00, 0x01, 0x02, 0x80:
			results[i] = SubscribeResult(d)
		default:
			return fmt.Errorf("invalid subscribe result: %d", d)
		}
	}
	p.Header.decode(b[0])
	p.MessageID = decodeMessageID(data[0:2])
	p.Results = results
	return nil
}

// AddResult adds a result of SUBSRIBE to the topic.
func (p *SubACK) AddResult(r SubscribeResult) {
	if p.Results == nil {
		p.Results = make([]SubscribeResult, 0, 4)
	}
	p.Results = append(p.Results, r)
}

// SubscribeResult represents result of subscribe to topic.
type SubscribeResult uint8

const (
	// SubscribeAtMostOnce is "Success - Maximum QoS 0"
	SubscribeAtMostOnce SubscribeResult = 0x00

	// SubscribeAtLeastOnce is "Success - Maximum QoS 1"
	SubscribeAtLeastOnce = 0x01

	// SubscribeExactOnce is "Success - Maximum QoS 2"
	SubscribeExactOnce = 0x02

	// SubscribeFailure is "Failure"
	SubscribeFailure = 0x80
)

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
	// TODO: rewrite with newDecoder()
	if len(b) < 2 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TUnsubscribe {
		return errors.New("type mismatch")
	}
	data, err := decodeRemain(b[2:])
	if err != nil {
		return nil
	}
	if len(data) < 2 {
		return errors.New("too short remain")
	}
	topics, err := decodeStrings(bytes.NewReader(data[2:]))
	if err != nil {
		return err
	}
	p.Header.decode(b[0])
	p.MessageID = decodeMessageID(data[0:2])
	p.Topics = topics
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
	d := newDecoder(b, TUnsubACK)
	packetID, err := d.readPacketID()
	if err != nil {
		return err
	}
	*p = UnsubACK{
		Header: d.header,
		MessageID: packetID,
	}
	return nil
}

// Topic represents topics to subscribe.
type Topic struct {
	Filter       string
	RequestedQoS QoS
}

func encodeTopics(topics []Topic) ([]byte, error) {
	buf := bytes.Buffer{}
	for i, t := range topics {
		n := encodeString(t.Filter)
		if n == nil {
			return nil, fmt.Errorf("too long topic name in #%d", i)
		}
		_, err := buf.Write(n)
		if err != nil {
			return nil, err
		}
		err = buf.WriteByte(byte(t.RequestedQoS & 0x03))
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func decodeTopic(r Reader) (*Topic, error) {
	s, err := decodeString(r)
	if err == io.EOF {
		return nil, nil
	}
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	return &Topic{
		Filter: s,
		RequestedQoS: QoS(b),
	}, nil
}

func decodeTopics(r Reader) ([]Topic, error) {
	var v []Topic
	for {
		t, err := decodeTopic(r)
		if err != nil {
			return nil, err
		}
		if t == nil {
			return v, nil
		}
		v = append(v, *t)
	}
}
