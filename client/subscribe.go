package client

// SubscribeError is detailed error for Subscribe().
type SubscribeError struct {
	// MismatchPacketID is set true, when detect mismatch of packet ID in ACK.
	MismatchPacketID    bool
	MismatchResultCount bool
	RequestedQoS        []QoS
	ResultQoS           []QoS
}

func (e *SubscribeError) Error() string {
	// FIXME: more detailed error message.
	return "something wrong on subscribe"
}

func (e *SubscribeError) hasErrors() bool {
	if e.MismatchPacketID {
		return true
	}
	if e.MismatchResultCount {
		return true
	}
	for i, q := range e.ResultQoS {
		if q == Failure || q != e.RequestedQoS[i] {
			return true
		}
	}
	return false
}

// UnsubscribeError is detailed error for Unsubscribe()
type UnsubscribeError struct {
	// MismatchPacketID is set true, when detect mismatch of packet ID in ACK.
	MismatchPacketID bool
}

func (e *UnsubscribeError) Error() string {
	if e.MismatchPacketID {
		return "mismatch packet ID"
	}
	return "unknown error"
}

func (e *UnsubscribeError) hasErrors() bool {
	if e.MismatchPacketID {
		return true
	}
	return false
}
