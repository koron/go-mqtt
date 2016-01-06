package packet

// Header represents common properties for all types of packet.
type Header struct {
	Type   Type
	Flags  uint8
	QoS    QoS
	Retain bool
}
