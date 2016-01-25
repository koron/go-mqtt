package server

// Topic represents a topic filter fanned in.
type Topic struct {

	// Filter is a filter string for topics which interested in.
	Filter string

	// QoS is required QoS for this topic filter.
	QoS QoS
}
