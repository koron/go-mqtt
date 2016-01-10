package packet

import "errors"

// Publish represents PUBLISH packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#publish
type Publish struct {
	Header
	TopicName string
	MessageID MessageID
	Payload   []byte
}

var _ Packet = (*Publish)(nil)

// Encode returns serialized Publish packet.
func (p *Publish) Encode() ([]byte, error) {
	var (
		header = &Header{
			Type:   p.Type,
			Dup:    p.Dup,
			QoS:    p.QoS,
			Retain: p.Retain,
		}
		topicName = encodeString(p.TopicName)
		messageID []byte
	)
	if topicName == nil {
		return nil, errors.New("too long TopicName")
	}
	if p.isMessageIDRequired(header.QoS) {
		messageID = p.MessageID.bytes()
	}
	return encode(header, topicName, messageID, p.Payload)
}

// Decode deserializes []byte as Publish packet.
func (p *Publish) Decode(b []byte) error {
	d, err := newDecoder(b, TPublish)
	if err != nil {
		return err
	}
	topicName, err := d.readString()
	if err != nil {
		return err
	}
	var messageID MessageID
	if p.isMessageIDRequired(d.header.QoS) {
		messageID, err = d.readPacketID()
		if err != nil {
			return err
		}
	}
	payload, err := d.readRemainBytes()
	if err != nil {
		return err
	}
	if err := d.finish(); err != nil {
		return err
	}
	*p = Publish{
		Header:    d.header,
		TopicName: topicName,
		MessageID: messageID,
		Payload:   payload,
	}
	return nil
}

func (p *Publish) isMessageIDRequired(qos QoS) bool {
	switch qos {
	case QAtLeastOnce, QExactlyOnce:
		return true
	default:
		return false
	}
}

// PubACK represents PUBACK packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#puback
type PubACK struct {
	Header
	MessageID MessageID
}

var _ Packet = (*PubACK)(nil)

// Encode returns serialized PubACK packet.
func (p *PubACK) Encode() ([]byte, error) {
	return encode(&Header{Type: TPubACK}, p.MessageID.bytes())
}

// Decode deserializes []byte as PubACK packet.
func (p *PubACK) Decode(b []byte) error {
	d, err := newDecoder(b, TPubACK)
	if err != nil {
		return err
	}
	packetID, err := d.readPacketID()
	if err != nil {
		return err
	}
	if err := d.finish(); err != nil {
		return err
	}
	*p = PubACK{
		Header:    d.header,
		MessageID: packetID,
	}
	return nil
}

// PubRec represents PUBREC packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#pubrec
type PubRec struct {
	Header
	MessageID MessageID
}

var _ Packet = (*PubRec)(nil)

// Encode returns serialized PubRec packet.
func (p *PubRec) Encode() ([]byte, error) {
	return encode(&Header{Type: TPubRec}, p.MessageID.bytes())
}

// Decode deserializes []byte as PubRec packet.
func (p *PubRec) Decode(b []byte) error {
	d, err := newDecoder(b, TPubRec)
	if err != nil {
		return err
	}
	packetID, err := d.readPacketID()
	if err != nil {
		return err
	}
	if err := d.finish(); err != nil {
		return err
	}
	*p = PubRec{
		Header:    d.header,
		MessageID: packetID,
	}
	return nil
}

// PubRel represents PUBREL packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#pubrel
type PubRel struct {
	Header
	MessageID MessageID
}

var _ Packet = (*PubRel)(nil)

// Encode returns serialized PubRel packet.
func (p *PubRel) Encode() ([]byte, error) {
	return encode(&Header{Type: TPubRel}, p.MessageID.bytes())
}

// Decode deserializes []byte as PubRel packet.
func (p *PubRel) Decode(b []byte) error {
	d, err := newDecoder(b, TPubRel)
	if err != nil {
		return err
	}
	packetID, err := d.readPacketID()
	if err != nil {
		return err
	}
	if err := d.finish(); err != nil {
		return err
	}
	*p = PubRel{
		Header:    d.header,
		MessageID: packetID,
	}
	return nil
}

// PubComp represents PUBCOMP packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#pubcomp
type PubComp struct {
	Header
	MessageID MessageID
}

var _ Packet = (*PubComp)(nil)

// Encode returns serialized PubComp packet.
func (p *PubComp) Encode() ([]byte, error) {
	return encode(&Header{Type: TPubComp}, p.MessageID.bytes())
}

// Decode deserializes []byte as PubComp packet.
func (p *PubComp) Decode(b []byte) error {
	d, err := newDecoder(b, TPubComp)
	if err != nil {
		return err
	}
	packetID, err := d.readPacketID()
	if err != nil {
		return err
	}
	if err := d.finish(); err != nil {
		return err
	}
	*p = PubComp{
		Header:    d.header,
		MessageID: packetID,
	}
	return nil
}
