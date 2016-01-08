package packet

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
