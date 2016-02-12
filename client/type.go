package client

// PublishedFunc is called when receive a message.
type PublishedFunc func(m *Message)

// DisconnectedFunc is called when a connection was lost.
// reason can be one of Reason or other errors.
type DisconnectedFunc func(reason error, param Param)

// Will represents MQTT's will message.
type Will struct {
	QoS     QoS
	Retain  bool
	Topic   string
	Message string
}

// Reason represents reason of disconnection.
type Reason int

const (
	// Explicitly shows called Disconnect() explicitly.
	Explicitly Reason = iota

	// Timeout shows by timeout.
	Timeout
)

func (r Reason) Error() string {
	switch r {
	case Explicitly:
		return "disconnected explicitly"
	case Timeout:
		return "detect timeout"
	default:
		return "unknown reason"
	}
}
