package packet

// QoS represents QoS levels of MQTT.
type QoS uint8

const (
	// QAtMostOnce represents "fire and forget" (<=1).
	QAtMostOnce QoS = iota

	// QAtLeastOnce represents "acknowledged delivery" (>=1).
	QAtLeastOnce

	// QExactlyOnce represents "assured delivery" (=1).
	QExactlyOnce

	// QReserved is reseved.
	QReserved
)

func (q QoS) String() string {
	switch q {
	case QAtMostOnce:
		return "[QoS 0 - at most once]"
	case QAtLeastOnce:
		return "[QoS 1 - at least once]"
	case QExactlyOnce:
		return "[QoS 2 - exactly once]"
	case QReserved:
		return "[QoS 3 - reserved]"
	default:
		return "[QoS UNDEFINED]"
	}
}
