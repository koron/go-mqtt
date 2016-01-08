package packet

import "fmt"

// Type is the type representing the MQTT packet's message type.
type Type uint8

const (
	// TReserved is a reserved value and should be considered an invalid message
	// type.
	TReserved Type = iota

	// TConnect is request to connect.  CtoS
	TConnect

	// TConnACK is ack for TConnect. StoC
	TConnACK

	// TPublish is publish a message. CtoS or StoC.
	TPublish

	// TPubACK is ack for TPublish for QoS1 message. CtoS or StoC.
	TPubACK

	// TPubRec is TPublish received for QoS2 messages. CtoS or StoC.
	TPubRec

	// TPubRel is TPublish release for QoS2 messages. CtoS or StoC.
	TPubRel

	// TPubComp is TPublish complete for QoS2 messages. CtoS or StoC.
	TPubComp

	// TSubscribe is value for subscribe a topic. CtoS
	TSubscribe

	// TSubACK is ack for TSubscribe. StoC
	TSubACK

	// TUnsubscribe is value for unsubscrbie a topic. CtoS.
	TUnsubscribe

	// TUnsubACK is ack for TUnsubscribe. StoC.
	TUnsubACK

	// TPingReq is value for PING response. CtoS.
	TPingReq

	// TPingResp is value for PING response. StoC.
	TPingResp

	// TDisconnect is value for disconnect.  CtoS.
	TDisconnect

	// TReserved2 is a reserved value and should be considered an invalid
	// message type.
	TReserved2
)

type typeDesc struct {
	Type
	Name  string
	Desc  string
	Flags uint8
}

var typeDescs = []*typeDesc{
	&typeDesc{
		Type:  TReserved,
		Name:  "RESERVED",
		Desc:  "Reserved",
		Flags: 0,
	},
	&typeDesc{
		Type:  TConnect,
		Name:  "CONNECT",
		Desc:  "Client request to connect to Server",
		Flags: 0,
	},
	&typeDesc{
		Type:  TConnACK,
		Name:  "CONNACK",
		Desc:  "Connect acknowledgement",
		Flags: 0,
	},
	&typeDesc{
		Type:  TPublish,
		Name:  "PUBLISH",
		Desc:  "Publish message",
		Flags: 0,
	},
	&typeDesc{
		Type:  TPubACK,
		Name:  "PUBACK",
		Desc:  "Publish acknowledgement",
		Flags: 0,
	},
	&typeDesc{
		Type:  TPubRec,
		Name:  "PUBREC",
		Desc:  "Publish received (assured delivery part 1)",
		Flags: 0,
	},
	&typeDesc{
		Type:  TPubRel,
		Name:  "PUBREL",
		Desc:  "Publish release (assured delivery part 2)",
		Flags: 2,
	},
	&typeDesc{
		Type:  TPubComp,
		Name:  "PUBCOMP",
		Desc:  "Publish complete (assured delivery part 3)",
		Flags: 0,
	},
	&typeDesc{
		Type:  TSubscribe,
		Name:  "SUBSCRIBE",
		Desc:  "Client subscribe request",
		Flags: 2,
	},
	&typeDesc{
		Type:  TSubACK,
		Name:  "SUBACK",
		Desc:  "Subscribe acknowledgement",
		Flags: 0,
	},
	&typeDesc{
		Type:  TUnsubscribe,
		Name:  "UNSUBSCRIBE",
		Desc:  "Unsubscribe request",
		Flags: 2,
	},
	&typeDesc{
		Type:  TUnsubACK,
		Name:  "UNSUBACK",
		Desc:  "Unsubscribe acknowledgement",
		Flags: 0,
	},
	&typeDesc{
		Type:  TPingReq,
		Name:  "PINGREQ",
		Desc:  "PING request",
		Flags: 0,
	},
	&typeDesc{
		Type:  TPingResp,
		Name:  "PINGRESP",
		Desc:  "PING response",
		Flags: 0,
	},
	&typeDesc{
		Type:  TDisconnect,
		Name:  "DISCONNECT",
		Desc:  "Client is disconnecting",
		Flags: 0,
	},
	&typeDesc{
		Type:  TReserved2,
		Name:  "RESERVED2",
		Desc:  "Reserved",
		Flags: 0,
	},
}

var typeUnknownDesc = &typeDesc{
	Type:  0,
	Name:  "UNKNOWN",
	Desc:  "UNKNOWN",
	Flags: 0,
}

func decodeType(b byte) Type {
	return Type(b >> 4 & 0x0f)
}

func (t Type) desc() *typeDesc {
	if t > TReserved2 {
		return typeUnknownDesc
	}
	return typeDescs[t]
}

// String returns string representation for a Type.
func (t Type) String() string {
	return t.desc().Name
}

// Name returns name of a Type.
func (t Type) Name() string {
	return t.desc().Name
}

// Flags returns the default flag values for message type.
func (t Type) Flags() uint8 {
	return t.desc().Flags
}

// NewPacket creates a new packet of this type.
func (t Type) NewPacket() (Packet, error) {
	switch t {
	case TConnect:
		return &Connect{Header: Header{Type: t}}, nil
	case TConnACK:
		return &ConnACK{Header: Header{Type: t}}, nil
	case TPublish:
		return &Publish{Header: Header{Type: t}}, nil
	case TPubACK:
		return &PubACK{Header: Header{Type: t}}, nil
	case TPubRec:
		return &PubRec{Header: Header{Type: t}}, nil
	case TPubRel:
		return &PubRel{Header: Header{Type: t}}, nil
	case TPubComp:
		return &PubComp{Header: Header{Type: t}}, nil
	case TSubscribe:
		return &Subscribe{Header: Header{Type: t}}, nil
	case TSubACK:
		return &SubACK{Header: Header{Type: t}}, nil
	case TUnsubscribe:
		return &Unsubscribe{Header: Header{Type: t}}, nil
	case TUnsubACK:
		return &UnsubACK{Header: Header{Type: t}}, nil
	case TPingReq:
		return &PingReq{Header: Header{Type: t}}, nil
	case TPingResp:
		return &PingResp{Header: Header{Type: t}}, nil
	case TDisconnect:
		return &Disconnect{Header: Header{Type: t}}, nil
	}
	return nil, fmt.Errorf("not defined type: %d", t)
}
