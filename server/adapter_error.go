package server

// AdapterError represents adapter error.
type AdapterError interface {
	Error() string

	// Continue returns true, if connection is continuable after this error.
	Continue() bool
}

type adapterError struct {
	message string
	cont    bool // true if this error is continuable
}

var _ AdapterError = (*adapterError)(nil)

var (
	// ErrDisconnected uses when the connection should be disconnected.
	ErrDisconnected AdapterError = &adapterError{
		message: "disconnected",
	}

	// ErrNotSuported uses when adapter not implemented yet.
	ErrNotSuported AdapterError = &adapterError{
		message: "not supported yet",
	}

	// ErrNotAcceptable uses when adapter couldn't accept the packet.
	ErrNotAcceptable AdapterError = &adapterError{
		message: "not acceptable packet",
	}
)

func (ae *adapterError) Error() string {
	return ae.message
}

func (ae *adapterError) Continue() bool {
	return ae.cont
}
