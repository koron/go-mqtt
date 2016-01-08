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

func remainReader(b []byte) (*bytes.Reader, error) {
	r := bytes.NewReader(b)
	u, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}
	if len(b)-r.Len() > 4 {
		return nil, errors.New("too long remain length")
	}
	if r.Len() != int(u) {
		return nil, errors.New("unmatch remain length")
	}
	return r, nil
}

// DEPRECATED
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
	ul, err := decodeUint16(r)
	if err == errInsufficientUint16 {
		return "", errInsufficientString
	} else if err != nil {
		return "", nil
	}
	l := int(ul)
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

var errInsufficientUint16 = errors.New("insufficient uint16")

func decodeUint16(r Reader) (uint16, error) {
	b1, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	b2, err := r.ReadByte()
	if err == io.EOF {
		return 0, errInsufficientUint16
	} else if err != nil {
		return 0, err
	}
	return uint16(b1)<<8 | uint16(b2), nil
}

func decodePacketID(r Reader) (MessageID, error) {
	u, err := decodeUint16(r)
	if err == errInsufficientUint16 {
		return 0, errors.New("insufficient MessageID")
	} else if err != nil {
		return 0, err
	}
	return MessageID(u), nil
}

func decodeHeader(b byte) (*Header, error) {
	return &Header{
		Type:   decodeType(b),
		Dup:    b&0x08 != 0,
		QoS:    QoS(b >> 1 & 0x3),
		Retain: b&0x01 != 0,
	}, nil
}

var (
	errInvalidPacketLength     = errors.New("invalid packet length")
	errTypeMismatch            = errors.New("type mismatch")
	errInsufficientRemainBytes = errors.New("insufficient remain bytes")
)

type decoder struct {
	header Header
	r      *bytes.Reader
	err    error
}

func newDecoder(b []byte, t Type) *decoder {
	if len(b) < 2 {
		return &decoder{err: errInvalidPacketLength}
	}
	h, err := decodeHeader(b[0])
	if err != nil {
		return &decoder{err: err}
	} else if h.Type != t {
		return &decoder{err: errTypeMismatch}
	}
	r, err := remainReader(b[1:])
	if err != nil {
		return &decoder{err: err}
	}
	return &decoder{header: *h, r: r}
}

func (d *decoder) readRemainBytes() ([]byte, error) {
	if d.err != nil {
		return nil, d.err
	}
	b := make([]byte, d.r.Len())
	n, err := d.r.Read(b)
	if err != nil {
		d.err = err
		return nil, err
	} else if n != len(b) {
		d.err = errInsufficientRemainBytes
		return nil, d.err
	}
	return b, nil
}

func (d *decoder) readPacketID() (MessageID, error) {
	if d.err != nil {
		return 0, d.err
	}
	id, err := decodePacketID(d.r)
	if err != nil {
		d.err = err
		return 0, d.err
	}
	return id, nil
}

func (d *decoder) readString() (string, error) {
	if d.err != nil {
		return "", d.err
	}
	s, err := decodeString(d.r)
	if err != nil {
		d.err = err
		return "", d.err
	}
	return s, nil
}
