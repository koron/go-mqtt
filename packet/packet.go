package packet

// Packet represents common I/F for packats.
type Packet interface {
	// Encode serializes packet to []byte.
	Encode() ([]byte, error)
}

// MessageID is identifier for packet/message.
type MessageID uint16

func (id MessageID) bytes() []byte {
	return []byte{
		byte(id >> 8 & 0xff),
		byte(id >> 0 & 0xff),
	}
}

