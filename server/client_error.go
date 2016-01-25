package server

type clientError struct {
	message string
	cont    bool // true if this error is continuable
}

var _ error = (*clientError)(nil)

var (
	errDisconnected = &clientError{
		message: "disconnected",
	}

	errNotSuported = &clientError{
		message: "not supported yet",
	}

	errNotAcceptable = &clientError{
		message: "not acceptable packet",
	}
)

func (ce *clientError) Error() string {
	return ce.message
}
