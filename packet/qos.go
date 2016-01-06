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
