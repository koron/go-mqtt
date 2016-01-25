package server

import "github.com/koron/go-mqtt/packet"

// PacketFilter filters all packets which recieve and send.
type PacketFilter interface {
	// PreProcess receives all packets after received and before it is
	// processed.
	PreProcess(p packet.Packet) error

	// PreSend is called before send a packet, can modify d:datagram.
	PreSend(p packet.Packet, d []byte) ([]byte, error)

	// PostSend is called after send a packet.
	PostSend(p packet.Packet, d []byte)
}

// ClientAdapter prorvides MQTT client adapter.
type ClientAdapter interface {

	// ID returns client ID.
	ID() string

	// IsSessionPresent returns true, if previous session is reverted.
	IsSessionPresent() bool

	// OnDisconnect is called when recieve DISCONNECT packet.
	OnDisconnect() error

	// OnPing is called when recieve PINGREQ packet.  If it returns false,
	// PINGRESP is not sent.
	OnPing() (bool, error)

	// OnSubscribe is called when receive SUBSCRIBE packet.
	OnSubscribe() error

	// OnUnsubscribe is called when receive UNSUBSCRIBE packet.
	OnUnsubscribe() error

	// OnPublish is called when receive PUBLISH packet.
	OnPublish() error
}

// NullClientAdapter is a default implementation of client adapter.
type NullClientAdapter struct {
	// ClientID holds client ID at Connect.
	ClientID string

	// SessionPresent indicates client having session info.
	SessionPresent bool
}

var (
	_ ClientAdapter = (*NullClientAdapter)(nil)
	_ PacketFilter = (*NullClientAdapter)(nil)
)

// ID returns client ID.
func (ca *NullClientAdapter) ID() string {
	return ca.ClientID
}

// IsSessionPresent returns true, if previous session is reverted.
func (ca *NullClientAdapter) IsSessionPresent() bool {
	return ca.SessionPresent
}

// PreProcess does nothing. just returns nil.
func (ca *NullClientAdapter) PreProcess(p packet.Packet) error {
	return nil
}

// PreSend does nothing.
func (ca *NullClientAdapter) PreSend(p packet.Packet, d []byte) ([]byte, error) {
	return d, nil
}

// PostSend does nothing.
func (ca *NullClientAdapter) PostSend(p packet.Packet, d []byte) {
}

// OnDisconnect does nothing.
func (ca *NullClientAdapter) OnDisconnect() error {
	return nil
}

// OnPing does nothing.
func (ca *NullClientAdapter) OnPing() (bool, error) {
	return true, nil
}

// OnSubscribe does nothing.
func (ca *NullClientAdapter) OnSubscribe() error {
	// TODO:
	return nil
}

// OnUnsubscribe does nothing.
func (ca *NullClientAdapter) OnUnsubscribe() error {
	// TODO:
	return nil
}

// OnPublish does nothing.
func (ca *NullClientAdapter) OnPublish() error {
	// TODO:
	return nil
}
