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
	ReturnCode ConnectReturnCode
}

var _ Packet = (*ConnACK)(nil)

// ConnectReturnCode is used in ConnACK. "Connect Return Code"
type ConnectReturnCode uint8

const (
	// ConnectAccept is "Connect Accepted".
	ConnectAccept ConnectReturnCode = iota

	// ConnectUnacceptableProtocolVersion is "Connection Refused: unacceptable protocol version"
	ConnectUnacceptableProtocolVersion

	// ConnectIdentifierRejected is "Connection Refused: identifier rejected"
	ConnectIdentifierRejected

	// ConnectServerUnavailable is "Connection Refused: server unavailable"
	ConnectServerUnavailable

	// ConnectBadUserNameOrPassword is "Connection Refused: bad user name or password"
	ConnectBadUserNameOrPassword

	// ConnectNotAuthorized is "Connection Refused: not authorized"
	ConnectNotAuthorized
)

// Encode returns serialized ConnACK packet.
func (p *ConnACK) Encode() ([]byte, error) {
	return encode(&Header{Type: TConnACK}, []byte{0x00, byte(p.ReturnCode)})
}

// Disconnect represents DISCONNECT packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#disconnect
type Disconnect struct {
	Header
}

var _ Packet = (*Disconnect)(nil)

// Encode returns serialized Disconnect packet.
func (p *Disconnect) Encode() ([]byte, error) {
	return encode(&Header{Type: TDisconnect}, nil)
}
