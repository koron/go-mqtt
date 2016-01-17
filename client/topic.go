package client

import "github.com/koron/go-mqtt/packet"

// Topic represents a topic filter fanned in.
type Topic struct {

	// Filter is a filter string for topics which interested in.
	Filter string

	// QoS is required QoS for this topic filter.
	QoS QoS
}

func (t *Topic) packetTopic() (packet.Topic, error) {
	return packet.Topic{
		Filter:       t.Filter,
		RequestedQoS: t.QoS.qos(),
	}, nil
}

func packetTopics(topics []Topic) ([]packet.Topic, error) {
	array := make([]packet.Topic, len(topics))
	for i, t := range topics {
		pt, err := t.packetTopic()
		if err != nil {
			return nil, err
		}
		array[i] = pt
	}
	return array, nil
}
