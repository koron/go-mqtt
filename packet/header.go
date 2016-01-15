package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

// header represents common properties for all types of packet.
type header struct {
	Type   Type
	Dup    bool
	QoS    QoS
	Retain bool
}

func encode(h *header, payloads ...[]byte) ([]byte, error) {
	buf := bytes.Buffer{}
	// encode byte1
	b := byte(h.Type)&0x0f<<4 + byte(h.QoS)&0x03<<1
	if h.Dup {
		b |= 0x08
	}
	if h.Retain {
		b |= 0x01
	}
	err := buf.WriteByte(b)
	if err != nil {
		return nil, err
	}
	// calculate remaining length.
	var rlen int
	for _, payload := range payloads {
		rlen += len(payload)
	}
	if rlen > math.MaxInt32 {
		return nil, errors.New("too long payload")
	} else if rlen == 0 {
		err = buf.WriteByte(0)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	// encode remaining length.
	b2 := make([]byte, binary.MaxVarintLen32)
	n := binary.PutUvarint(b2, uint64(rlen))
	_, err = buf.Write(b2[:n])
	if err != nil {
		return nil, err
	}
	// encode all payloads.
	for _, payload := range payloads {
		if payload == nil {
			continue
		}
		_, err = buf.Write(payload)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// decode decodes a byte as header except message type.
func (h *header) decode(b byte) {
	h.Dup = b&0x08 != 0
	h.QoS = QoS(b >> 1 & 0x3)
	h.Retain = b&0x01 != 0
}
