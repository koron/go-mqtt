package server

import "github.com/surgemq/message"

// ConnectHandler checks a connection is acceptable or not.
type ConnectHandler func(srv *Server, conn PreConn, msg *message.ConnectMessage) error

// ReceiveHandler called when receive a MQTT message.
type ReceiveHandler func(conn Conn, msg message.Message) error
