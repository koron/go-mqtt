package client

// Message represents a MQTT's published message.
type Message struct {
	Topic string
	Body  []byte
}
