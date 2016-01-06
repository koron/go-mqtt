package packet

// Subscribe represents SUBSRIBE packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#subscribe
type Subscribe struct {
	Header
	// TODO: add props for Subscribe.
}

var _ Packet = (*Subscribe)(nil)

// Encode returns serialized Subscribe packet.
func (p *Subscribe) Encode() ([]byte, error) {
	return nil, nil
}

// SubACK represents SUBACK packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#suback
type SubACK struct {
	Header
	// TODO: add props for SubACK.
}

var _ Packet = (*SubACK)(nil)

// Encode returns serialized SubACK packet.
func (p *SubACK) Encode() ([]byte, error) {
	return nil, nil
}

// Unsubscribe represents UNSUBSCRIBE packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#unsubscribe
type Unsubscribe struct {
	Header
	// TODO: add props for Unsubscribe.
}

var _ Packet = (*Unsubscribe)(nil)

// Encode returns serialized Unsubscribe packet.
func (p *Unsubscribe) Encode() ([]byte, error) {
	return nil, nil
}

// UnsubACK represents UNSUBACK packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#unsuback
type UnsubACK struct {
	Header
	// TODO: add props for UnsubACK.
}

var _ Packet = (*UnsubACK)(nil)

// Encode returns serialized UnsubACK packet.
func (p *UnsubACK) Encode() ([]byte, error) {
	return nil, nil
}
