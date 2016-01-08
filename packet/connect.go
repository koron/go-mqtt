package packet

import (
	"errors"
	"fmt"
)

const protocolName = "MQIsdp"

// Connect represents CONNECT packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#connect
type Connect struct {
	Header
	ClientID     string
	Version      uint8
	Username     *string
	Password     *string
	CleanSession bool
	KeepAlive    uint16
	WillFlag     bool
	WillQoS      QoS
	WillRetain   bool
	WillTopic    string
	WillMessage  string
}

var _ Packet = (*Connect)(nil)

// Encode returns serialized Connect packet.
func (p *Connect) Encode() ([]byte, error) {
	var (
		header       = &Header{Type: TConnect}
		clientID     = encodeString(p.ClientID)
		willTopic    []byte
		willMessage  []byte
		username     []byte
		password     []byte
		connectFlags byte
	)
	if l := len(p.ClientID); l <= 0 || l > 23 {
		return nil, errors.New("too short/long ClientID")
	}
	if p.Username != nil {
		username = encodeString(*p.Username)
		if username == nil {
			return nil, errors.New("too long Username")
		}
		connectFlags |= 0x80
	}
	if p.Password != nil {
		password = encodeString(*p.Password)
		if password == nil {
			return nil, errors.New("too long Password")
		}
		connectFlags |= 0x40
	}
	if p.WillFlag {
		willTopic = encodeString(p.WillTopic)
		if willTopic == nil {
			return nil, errors.New("too long WillTopic")
		}
		willMessage = encodeString(p.WillMessage)
		if willMessage == nil {
			return nil, errors.New("too long WillMessage")
		}
		connectFlags |= (byte)(p.WillQoS&0x03<<3) | 0x04
		if p.WillRetain {
			connectFlags |= 0x20
		}
	}
	if p.CleanSession {
		connectFlags |= 0x02
	}
	return encode(
		header,
		encodeString(protocolName),
		[]byte{byte(p.Version), connectFlags},
		encodeUint16(p.KeepAlive),
		clientID,
		willTopic,
		willMessage,
		username,
		password)
}

// Decode deserializes []byte as Connect packet.
func (p *Connect) Decode(b []byte) error {
	// TODO: implement me.
	return nil
}

// ConnACK represents CONNACK packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#connack
type ConnACK struct {
	Header
	SessionPresent bool
	ReturnCode     ConnectReturnCode
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
	var flags byte
	if p.SessionPresent {
		flags |= 0x01
	}
	return encode(&Header{Type: TConnACK}, []byte{flags, byte(p.ReturnCode)})
}

// Decode deserializes []byte as ConnACK packet.
func (p *ConnACK) Decode(b []byte) error {
	// TODO: rewrite with newDecoder()
	if len(b) != 4 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TConnACK {
		return errors.New("type mismatch")
	}
	if b[1] != 2 {
		return errors.New("invalid remain length")
	}
	if c := ConnectReturnCode(b[3]); c > ConnectNotAuthorized {
		return fmt.Errorf("invalid return code: %d", c)
	}
	p.Header.decode(b[0])
	p.SessionPresent = b[2] & 0x01 != 0
	p.ReturnCode = ConnectReturnCode(b[3])
	return nil
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

// Decode deserializes []byte as Disconnect packet.
func (p *Disconnect) Decode(b []byte) error {
	// TODO: rewrite with newDecoder()
	if len(b) != 2 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TDisconnect {
		return errors.New("type mismatch")
	}
	if b[1] != 0 {
		return errors.New("invalid remain length")
	}
	p.Header.decode(b[0])
	return nil
}
