package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Reader declares stream which can be decoded as Packets.
// bufio.Reader is one of satisfied implementations.
type Reader interface {
	io.Reader
	io.ByteReader
}

// SplitDecode splits datagram from Reader and decode it as a Packet.
func SplitDecode(r Reader) (Packet, error) {
	b, err := Split(r)
	if err != nil {
		return nil, err
	}
	return Decode(b)
}

// Split splits datagram of a Packet from Reader.
func Split(r Reader) ([]byte, error) {
	buf := &bytes.Buffer{}
	// read header: message type
	c, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	err = buf.WriteByte(c)
	if err != nil {
		return nil, err
	}
	// read length of payload
	l, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}
	// FIXME: should be check l is too long.
	b := make([]byte, binary.MaxVarintLen32)
	n := binary.PutUvarint(b, l)
	buf.Write(b[0:n])
	// read whole payload to buf.
	_, err = io.CopyN(buf, r, int64(l))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode decodes a Packet from datagram.
func Decode(b []byte) (Packet, error) {
	if len(b) < 2 {
		return nil, errors.New("too short []byte")
	}
	t := decodeType(b[0])
	p, err := t.NewPacket()
	if err != nil {
		return nil, err
	}
	err = p.Decode(b)
	if err != nil {
		return nil, err
	}
	return p, nil
}
