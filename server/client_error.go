package server

type clientError struct {
	message string
	cont    bool // true if this error is continuable
}

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
var _ error = (*clientError)(nil)

func (ce *clientError) Error() string {
	return ce.message
}
