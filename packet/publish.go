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
	d := newDecoder(b, TPublish)
	var (
		topicName string
		messageID MessageID
		payload   []byte
	)
	topicName, _ = d.readString()
	if p.isMessageIDRequired(d.header.QoS) {
		messageID, _ = d.readPacketID()
	}
	payload, err := d.readRemainBytes()
	if err != nil {
		if err == errInsufficientRemainBytes {
			err = errors.New("insufficient payload")
		}
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
	if len(b) != 4 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TPubACK {
		return errors.New("type mismatch")
	}
	if b[1] != 2 {
		return errors.New("invalid remain length")
	}
	p.Header.decode(b[0])
	p.MessageID = decodeMessageID(b[2:])
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
	if len(b) != 4 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TPubRec {
		return errors.New("type mismatch")
	}
	if b[1] != 2 {
		return errors.New("invalid remain length")
	}
	p.Header.decode(b[0])
	p.MessageID = decodeMessageID(b[2:])
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
	if len(b) != 4 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TPubRel {
		return errors.New("type mismatch")
	}
	if b[1] != 2 {
		return errors.New("invalid remain length")
	}
	p.Header.decode(b[0])
	p.MessageID = decodeMessageID(b[2:])
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
	if len(b) != 4 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TPubComp {
		return errors.New("type mismatch")
	}
	if b[1] != 2 {
		return errors.New("invalid remain length")
	}
	p.Header.decode(b[0])
	p.MessageID = decodeMessageID(b[2:])
	return nil
}
