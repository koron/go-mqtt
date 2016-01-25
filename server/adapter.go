package server

import "github.com/koron/go-mqtt/packet"

// Adapter provides MQTT server (broker) adapter.
type Adapter interface {

	// Connect is called when a new client try to connect MQTT broker.
	// It can return one of ConnectError.
	// ClientAdapter can implement PacketFilter.
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
