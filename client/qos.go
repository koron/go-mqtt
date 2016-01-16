package client

import "github.com/koron/go-mqtt/packet"

// QoS represents QoS levels of MQTT.
// Invalid values are treated as AtMostOnce.
type QoS uint8

const (
	// AtMostOnce represents "fire and forget" (<=1).
	AtMostOnce QoS = iota

	// AtLeastOnce represents "acknowledged delivery" (>=1).
	AtLeastOnce

	// ExactlyOnce represents "assured delivery" (=1).
	ExactlyOnce
)

func (q QoS) String() string {
	switch q {
	case AtLeastOnce:
		return "at least once"
	case AtMostOnce:
		return "at most once"
	case ExactlyOnce:
		return "exactly once"
	default:
		return "at least once (fallback)"
	}
}

func (q QoS) qos() packet.QoS {
	switch q {
	case AtLeastOnce:
		return packet.QAtLeastOnce
	case AtMostOnce:
		return packet.QAtMostOnce
	case ExactlyOnce:
		return packet.QExactlyOnce
	default:
		return packet.QAtLeastOnce
	}
}
