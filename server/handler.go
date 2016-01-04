package server

import "github.com/surgemq/message"

// ConnectHandler checks a connection is acceptable or not.
type ConnectHandler func(srv *Server, conn Conn, msg *message.ConnectMessage) error
