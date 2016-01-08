package packet

import "errors"

var (
	// ErrTooLongPayload shows payload is too long to encode.
	ErrTooLongPayload = errors.New("too long payload")
)
