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

func decodeRemain(b []byte) ([]byte, error) {
	r := bytes.NewReader(b)
	u, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}
	if r.Len() > 4 {
		return nil, errors.New("too long remain length")
	}
	if len(b)-r.Len() != int(u) {
		return nil, errors.New("unmatch remain length")
	}
	return b[r.Len():], nil
}

var errInsufficientString = errors.New("insufficient string")

func decodeString(r Reader) (string, error) {
	b1, err := r.ReadByte()
	if err != nil {
		return "", err
	}
	b2, err := r.ReadByte()
	if err != nil {
		if err == io.EOF {
			err = errInsufficientString
		}
		return "", err
	}
	l := int(b1)<<8 | int(b2)
	b := make([]byte, l)
	n, err := r.Read(b)
	if err != nil {
		if err == io.EOF {
			err = errInsufficientString
		}
		return "", err
	}
	if n != l {
		return "", errInsufficientString
	}
	return string(b), nil
}

func decodeStrings(r Reader) ([]string, error) {
	var v []string
	for {
		s, err := decodeString(r)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		v = append(v, s)
	}
	return v, nil
}
