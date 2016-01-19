package server

import "github.com/koron/go-mqtt/packet"

// ConnectError represents connection error.
// It would be returned by Adapter#Connect().
type ConnectError int

var _ error = (ConnectError)(0)

const (
	// ErrUnacceptableProtocolVersion is "Connection Refused: unacceptable protocol version"
	ErrUnacceptableProtocolVersion ConnectError = iota + 1

	// ErrIdentifierRejected is "Connection Refused: identifier rejected"
	ErrIdentifierRejected

	// ErrServerUnavailable is "Connection Refused: server unavailable"
	ErrServerUnavailable

	// ErrBadUserNameOrPassword is "Connection Refused: bad user name or password"
	ErrBadUserNameOrPassword

	// ErrNotAuthorized is "Connection Refused: not authorized"
	ErrNotAuthorized
)

func (ce ConnectError) Error() string {
	switch ce {
	case ErrUnacceptableProtocolVersion:
		return "unacceptable protocol version"
	case ErrIdentifierRejected:
		return "identifier rejected"
	case ErrServerUnavailable:
		return "server unavailable"
	case ErrBadUserNameOrPassword:
		return "bad user name or password"
	case ErrNotAuthorized:
		return "not authorized"
	default:
		return "unknown"
	}
}

func (ce ConnectError) toRC() packet.ConnectReturnCode {
	switch ce {
	case ErrUnacceptableProtocolVersion:
		return packet.ConnectUnacceptableProtocolVersion
	case ErrIdentifierRejected:
		return packet.ConnectIdentifierRejected
	case ErrServerUnavailable:
		return packet.ConnectServerUnavailable
	case ErrBadUserNameOrPassword:
		return packet.ConnectBadUserNameOrPassword
	case ErrNotAuthorized:
		return packet.ConnectNotAuthorized
	default:
		return packet.ConnectNotAuthorized
	}
}
