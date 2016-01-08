package packet

import "math"

// Packet represents common I/F for packats.
type Packet interface {
	// Encode serializes packet to []byte.
	Encode() ([]byte, error)

	// Decode deserializes []byte as an packet.
	Decode([]byte) error
}

// MessageID is identifier for packet/message.
type MessageID uint16

func (id MessageID) bytes() []byte {
	return []byte{
		byte(id >> 8 & 0xff),
		byte(id >> 0 & 0xff),
	}
}

func encodeUint16(n uint16) []byte {
	return []byte{
		byte(n >> 8 & 0xff),
		byte(n >> 0 & 0xff),
	}
}

func encodeString(s string) []byte {
	l := len(s)
	if len(s) > math.MaxUint16 {
		return nil
	}
	b := make([]byte, l+2)
	b[0] = byte(l >> 8 & 0xff)
	b[1] = byte(l >> 0 & 0xff)
	copy(b[2:], []byte(s))
	return b
}
