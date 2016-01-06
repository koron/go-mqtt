package packet

// Connect represents CONNECT packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#connect
type Connect struct {
	Header
	// TODO: add props for Connect.
}

var _ Packet = (*Connect)(nil)

// Encode returns serialized Connect packet.
func (p *Connect) Encode() ([]byte, error) {
	return nil, nil
}

// ConnACK represents CONNACK packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#connack
type ConnACK struct {
	Header
	// TODO: add props for ConnACK.
}

var _ Packet = (*ConnACK)(nil)

// Encode returns serialized ConnACK packet.
func (p *ConnACK) Encode() ([]byte, error) {
	return nil, nil
}

// Disconnect represents DISCONNECT packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#disconnect
type Disconnect struct {
	Header
	// TODO: add props for Disconnect.
}

var _ Packet = (*Disconnect)(nil)

// Encode returns serialized Disconnect packet.
func (p *Disconnect) Encode() ([]byte, error) {
	return nil, nil
}
