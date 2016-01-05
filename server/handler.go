package server

import "github.com/surgemq/message"

// ConnectHandler checks a connection is acceptable or not.
type ConnectHandler func(srv *Server, conn PreConn, msg *message.ConnectMessage) error

// DisconnectedHandler notifies a disconnection.  msg can be nil when
// disconnected without DISCONNECT message.
type DisconnectedHandler func(srv *Server, conn DisConn, msg *message.DisconnectMessage) error

// ReceiveHandler called when receive a MQTT message.
type ReceiveHandler func(conn Conn, msg message.Message) error

// SentHandler called after sent a MQTT message.
type SentHandler func(conn Conn, msg message.Message) error

// SubscribleHandler called each topic to subscribe.
type SubscribleHandler func(conn Conn, topic string, requestedQos byte) (qos byte, err error)
