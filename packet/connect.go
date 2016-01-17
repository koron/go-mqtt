package packet

import (
	"errors"
	"fmt"
)

const (
	protocolVersion3 = 3
	protocolName3    = "MQIsdp"
	protocolVersion4 = 4
	protocolName4    = "MQTT"
)

// Connect represents CONNECT packet.
type Connect struct {
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
		header       = &header{Type: TConnect}
		protocolName string
		clientID     = encodeString(p.ClientID)
		willTopic    []byte
		willMessage  []byte
		username     []byte
		password     []byte
		connectFlags byte
	)
	switch p.Version {
	case protocolVersion3:
		protocolName = protocolName3
	case protocolVersion4:
		protocolName = protocolName4
	default:
		return nil, errors.New("unsupported protocol version")
	}
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
	d, err := newDecoder(b, TConnect)
	if err != nil {
		return err
	}
	protocolName, err := d.readString()
	if err != nil {
		return err
	}
	version, err := d.readByte()
	if err != nil {
		return err
	}
	switch version {
	case protocolVersion3:
		if protocolName != protocolName3 {
			return errors.New("mismatch protocol name and version")
		}
	case protocolVersion4:
		if protocolName != protocolName4 {
			return errors.New("mismatch protocol name and version")
		}
	default:
		return errors.New("unsupported protocol version")
	}
	connectFlags, err := d.readByte()
	if err != nil {
		return err
	}
	var (
		usernameFlag = connectFlags&0x80 != 0
		passwordFlag = connectFlags&0x40 != 0
		willRetain   = connectFlags&0x20 != 0
		willQoS      = QoS(connectFlags & 0x18 >> 3)
		willFlag     = connectFlags&0x04 != 0
		cleanSession = connectFlags&0x02 != 0
	)
	keepAlive, err := d.readUint16()
	if err != nil {
		return err
	}
	clientID, err := d.readString()
	if err != nil {
		return err
	}
	var (
		willTopic   string
		willMessage string
		username    *string
		password    *string
	)
	if willFlag {
		willTopic, err = d.readString()
		if err != nil {
			return err
		}
		willMessage, err = d.readString()
		if err != nil {
			return err
		}
	}
	if usernameFlag {
		s, err := d.readString()
		if err != nil {
			return err
		}
		username = &s
	}
	if passwordFlag {
		s, err := d.readString()
		if err != nil {
			return err
		}
		password = &s
	}
	if err := d.finish(); err != nil {
		return err
	}
	*p = Connect{
		ClientID:     clientID,
		Version:      uint8(version),
		Username:     username,
		Password:     password,
		CleanSession: cleanSession,
		KeepAlive:    keepAlive,
		WillFlag:     willFlag,
		WillQoS:      willQoS,
		WillRetain:   willRetain,
		WillTopic:    willTopic,
		WillMessage:  willMessage,
	}
	return nil
}

// ConnACK represents CONNACK packet.
type ConnACK struct {
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

func (c ConnectReturnCode) Error() string {
	switch c {
	case ConnectAccept:
		return "accepted"
	case ConnectUnacceptableProtocolVersion:
		return "unacceptable protocol version"
	case ConnectIdentifierRejected:
		return "identifier rejected"
	case ConnectServerUnavailable:
		return "server unavailable"
	case ConnectBadUserNameOrPassword:
		return "bad username or password"
	case ConnectNotAuthorized:
		return "not authorized"
	default:
		return "unknown connect return code"
	}
}

// Encode returns serialized ConnACK packet.
func (p *ConnACK) Encode() ([]byte, error) {
	var flags byte
	if p.SessionPresent {
		flags |= 0x01
	}
	return encode(&header{Type: TConnACK}, []byte{flags, byte(p.ReturnCode)})
}

// Decode deserializes []byte as ConnACK packet.
func (p *ConnACK) Decode(b []byte) error {
	d, err := newDecoder(b, TConnACK)
	if err != nil {
		return err
	}
	f, err := d.readByte()
	if err != nil {
		return err
	}
	sessionPresent := f&0x01 != 0
	c, err := d.readByte()
	if err != nil {
		return err
	}
	returnCode := ConnectReturnCode(c)
	if returnCode > ConnectNotAuthorized {
		return fmt.Errorf("invalid return code: %d", c)
	}
	if err := d.finish(); err != nil {
		return err
	}
	*p = ConnACK{
		SessionPresent: sessionPresent,
		ReturnCode:     returnCode,
	}
	return nil
}

// Disconnect represents DISCONNECT packet.
type Disconnect struct {
}

var _ Packet = (*Disconnect)(nil)

// Encode returns serialized Disconnect packet.
func (p *Disconnect) Encode() ([]byte, error) {
	return encode(&header{Type: TDisconnect}, nil)
}

// Decode deserializes []byte as Disconnect packet.
func (p *Disconnect) Decode(b []byte) error {
	d, err := newDecoder(b, TDisconnect)
	if err != nil {
		return err
	}
	if err := d.finish(); err != nil {
		return err
	}
	*p = Disconnect{}
	return nil
}
