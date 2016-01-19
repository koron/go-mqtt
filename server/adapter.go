package server

import "github.com/koron/go-mqtt/packet"

// Adapter provides MQTT server (broker) adapter.
type Adapter interface {

	// Connect is called when a new client try to connect MQTT broker.
	// It can return one of ConnectError.
	Connect(srv *Server, p *packet.Connect) (ClientAdapter, error)

	// Disconnect is called when a client disconnected.
	Disconnect(srv *Server, ca ClientAdapter, err error)
}

// NullAdapter is a default implementation of server adapter.
type NullAdapter struct {
}

var _ Adapter = (*NullAdapter)(nil)

// Connect is called when a new client try to connect MQTT broker.
func (a *NullAdapter) Connect(srv *Server, p *packet.Connect) (ClientAdapter, error) {
	return &NullClientAdapter{}, nil
}

// Disconnect is called when a client disconnected.
func (a *NullAdapter) Disconnect(srv *Server, ca ClientAdapter, err error) {
}

// DefaultAdapter is a default server adapter.
var DefaultAdapter = &NullAdapter{}

// ClientAdapter prorvides MQTT client adapter.
type ClientAdapter interface {

	// ID returns client ID.
	ID() string

	// IsSessionPresent() returns true, if previous session is reverted.
	IsSessionPresent() bool
}

// NullClientAdapter is a default implementation of client adapter.
type NullClientAdapter struct {
	// ClientID holds client ID at Connect.
	ClientID string

	// SessionPresent indicates client having session info.
	SessionPresent bool
}

var _ ClientAdapter = (*NullClientAdapter)(nil)

// ID returns client ID.
func (ca *NullClientAdapter) ID() string {
	return ca.ClientID
}

// IsSessionPresent returns true, if previous session is reverted.
func (ca *NullClientAdapter) IsSessionPresent() bool {
	return ca.SessionPresent
}
