package packet

import "io"

// Split splits datagram of a Packet from io.Reader.
func Split(r io.Reader) ([]byte, error) {
	// TODO:
	return nil, nil
}

// Decode decodes a Packet from datagram.
func Decode(b []byte) (Packet, error) {
	// TODO:
	return nil, nil
}

// SplitDecode splits datagram from io.Reader and decode it as a Packet.
func SplitDecode(r io.Reader) (Packet, error) {
	b, err := Split(r)
	if err != nil {
		return nil, err
	}
	return Decode(b)
}
