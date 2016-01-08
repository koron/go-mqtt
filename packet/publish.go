package packet

// Publish represents PUBLISH packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#publish
type Publish struct {
	Header
	// TODO: add props for Publish.
}

var _ Packet = (*Publish)(nil)

// Encode returns serialized Publish packet.
func (p *Publish) Encode() ([]byte, error) {
	return nil, nil
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
