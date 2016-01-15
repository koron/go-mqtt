package packet

// PingReq represents PINGREQ packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#pingreq
type PingReq struct {
}

var _ Packet = (*PingReq)(nil)

// Encode returns serialized PingReq packet.
func (p *PingReq) Encode() ([]byte, error) {
	return encode(&Header{Type: TPingReq}, nil)
}

// Decode deserializes []byte as PingReq packet.
func (p *PingReq) Decode(b []byte) error {
	d, err := newDecoder(b, TPingReq)
	if err != nil {
		return err
	}
	if err := d.finish(); err != nil {
		return err
	}
	*p = PingReq{}
	return nil
}

// PingResp represents PINGRESP packet.
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#pingresp
type PingResp struct {
}

var _ Packet = (*PingResp)(nil)

// Encode returns serialized PingResp packet.
func (p *PingResp) Encode() ([]byte, error) {
	return encode(&Header{Type: TPingResp}, nil)
}

// Decode deserializes []byte as PingResp packet.
func (p *PingResp) Decode(b []byte) error {
	d, err := newDecoder(b, TPingResp)
	if err != nil {
		return err
	}
	if err := d.finish(); err != nil {
		return err
	}
	*p = PingResp{}
	return nil
}
