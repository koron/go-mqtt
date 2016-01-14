package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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

var (
	errInvalidPacketLength     = errors.New("invalid packet length")
	errTypeMismatch            = errors.New("type mismatch")
	errInvalidRemainLength     = errors.New("invalid remain length")
	errInsufficientRemainBytes = errors.New("insufficient remain bytes")
	errInsufficientUint16      = errors.New("insufficient uint16")
	errInsufficientString      = errors.New("insufficient string")
	errInsufficientPacketID    = errors.New("insufficient PacketID")
	errUnreadBytes             = errors.New("unread bytes")
)

type decoder struct {
	header Header
	r      *bytes.Reader
}

func newDecoder(b []byte, t Type) (*decoder, error) {
	if len(b) < 2 {
		return nil, errInvalidPacketLength
	}
	b0 := b[0]
	h := Header{
		Type:   decodeType(b0),
		Dup:    b0&0x08 != 0,
		QoS:    QoS(b0 >> 1 & 0x3),
		Retain: b0&0x01 != 0,
	}
	if h.Type != t {
		return nil, errTypeMismatch
	}
	b = b[1:]
	r := bytes.NewReader(b)
	u, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}
	if len(b)-r.Len() > 4 || r.Len() != int(u) {
		return nil, errInvalidRemainLength
	}
	return &decoder{header: h, r: r}, nil
}

func (d *decoder) remainLen() int {
	return d.r.Len()
}

func (d *decoder) readByte() (byte, error) {
	b, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}
	return b, nil
}

func (d *decoder) readUint16() (uint16, error) {
	b1, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}
	b2, err := d.r.ReadByte()
	if err != nil {
		if err == io.EOF {
			err = errInsufficientUint16
		}
		return 0, err
	}
	return uint16(b1)<<8 | uint16(b2), nil
}

func (d *decoder) readPacketID() (ID, error) {
	id, err := d.readUint16()
	if err != nil {
		if err == errInsufficientUint16 {
			err = errInsufficientPacketID
		}
		return 0, err
	}
	return ID(id), nil
}

func (d *decoder) readString() (string, error) {
	ul, err := d.readUint16()
	if err != nil {
		if err == errInsufficientUint16 {
			err = errInsufficientString
		}
		return "", err
	}
	l := int(ul)
	b := make([]byte, l)
	n, err := d.r.Read(b)
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

func (d *decoder) readStrings() ([]string, error) {
	var v []string
	for {
		s, err := d.readString()
		if err == io.EOF {
			return v, nil
		} else if err != nil {
			return nil, err
		}
		v = append(v, s)
	}
}

func (d *decoder) readSubscribeResults() ([]SubscribeResult, error) {
	results := make([]SubscribeResult, 0, d.remainLen())
	for {
		b, err := d.readByte()
		if err == io.EOF {
			return results, nil
		} else if err != nil {
			return nil, err
		}
		switch b {
		case 0x00, 0x01, 0x02, 0x80:
			results = append(results, SubscribeResult(b))
		default:
			return nil, fmt.Errorf("invalid subscribe result: %d", b)
		}
	}
}

func (d *decoder) readRemainBytes() ([]byte, error) {
	b := make([]byte, d.r.Len())
	n, err := d.r.Read(b)
	if err != nil {
		return nil, err
	} else if n != len(b) {
		return nil, errInsufficientRemainBytes
	}
	return b, nil
}

func (d *decoder) readTopic() (*Topic, error) {
	s, err := d.readString()
	if err == io.EOF {
		return nil, nil
	}
	b, err := d.readByte()
	if err != nil {
		return nil, err
	}
	return &Topic{
		Filter:       s,
		RequestedQoS: QoS(b),
	}, nil
}

func (d *decoder) readTopics() ([]Topic, error) {
	var v []Topic
	for {
		t, err := d.readTopic()
		if err != nil {
			return nil, err
		}
		if t == nil {
			return v, nil
		}
		v = append(v, *t)
	}
}

func (d *decoder) finish() error {
	if d.r.Len() > 0 {
		return errUnreadBytes
	}
	return nil
}
