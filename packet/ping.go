package packet

import "errors"

// PingReq represents PINGREQ packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#pingreq
type PingReq struct {
	Header
}

var _ Packet = (*PingReq)(nil)

// Encode returns serialized PingReq packet.
func (p *PingReq) Encode() ([]byte, error) {
	return encode(&Header{Type: TPingReq}, nil)
}

// Decode deserializes []byte as PingReq packet.
func (p *PingReq) Decode(b []byte) error {
	if len(b) != 2 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TPingReq {
		return errors.New("type mismatch")
	}
	if b[1] != 0 {
		return errors.New("invalid remain length")
	}
	p.Header.decode(b[0])
	return nil
}

// PingResp represents PINGRESP packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#pingresp
type PingResp struct {
	Header
}

var _ Packet = (*PingResp)(nil)

// Encode returns serialized PingResp packet.
func (p *PingResp) Encode() ([]byte, error) {
	return encode(&Header{Type: TPingResp}, nil)
}

// Decode deserializes []byte as PingResp packet.
func (p *PingResp) Decode(b []byte) error {
	if len(b) != 2 {
		return errors.New("invalid packet length")
	}
	if decodeType(b[0]) != TPingResp {
		return errors.New("type mismatch")
	}
	if b[1] != 0 {
		return errors.New("invalid remain length")
	}
	p.Header.decode(b[0])
	return nil
}
