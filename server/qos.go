package server

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

	// Failure indicates "subscription failed".
	Failure = 0x80
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

func (q QoS) toSubscribeResult() packet.SubscribeResult {
	switch q {
	case AtMostOnce:
		return packet.SubscribeAtMostOnce
	case AtLeastOnce:
		return packet.SubscribeAtLeastOnce
	case ExactlyOnce:
		return packet.SubscribeExactOnce
	default:
		return packet.SubscribeFailure
	}
}

func toQoS(v packet.QoS) QoS {
	switch v {
	case packet.QAtMostOnce:
		return AtMostOnce
	case packet.QAtLeastOnce:
		return AtLeastOnce
	case packet.QExactlyOnce:
		return ExactlyOnce
	default:
		return Failure
	}
}
