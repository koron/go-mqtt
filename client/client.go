package client

import "github.com/koron/go-mqtt/packet"

// Client is a MQTT client.
type Client interface {
	// Disconnect shutdown MQTT connection.
	Disconnect(force bool) error

	// Ping sends a PING packet.
	Ping() error

	// Subscribe subsribes to topics.
	Subscribe(topics []Topic) error

	// Unsubscribe unsubscribes from topics.
	Unsubscribe(topics []string) error

	// Publish publishes a message to MQTT broker.
	Publish(qos QoS, retain bool, topic string, msg []byte) error

	// ReadMessage returns a message if it was available.  Otherwise this will
	// block.
	ReadMessage() (*Message, error)

	// PeekMessage returns true if ReadMessage() can return one or more
	// messages without blocking.
	PeekMessage() bool
}

// client implements a simple MQTT client.
type client struct {
	// TODO:
	nextID packet.ID
}

var _ Client = (*client)(nil)

func (c *client) Disconnect(force bool) error {
	// TODO:
	return nil
}

func (c *client) Ping() error {
	// TODO:
	return nil
}

func (c *client) Subscribe(topics []Topic) error {
	// TODO:
	return nil
}

func (c *client) Unsubscribe(topics []string) error {
	// TODO:
	return nil
}

func (c *client) Publish(qos QoS, retain bool, topic string, msg []byte) error {
	// TODO:
	return nil
}

func (c *client) PeekMessage() bool {
	// TODO:
	return false
}

func (c *client) ReadMessage() (*Message, error) {
	// TODO:
	return nil, nil
}
