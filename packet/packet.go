package packet

// Packet represents common I/F for packats.
type Packet interface {
	// Encode serializes packet to []byte.
	Encode() ([]byte, error)
}
