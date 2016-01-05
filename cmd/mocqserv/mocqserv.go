package main

import (
	"log"

	"github.com/koron/go-mqtt/server"
	"github.com/surgemq/message"
)

func recv(conn server.Conn, msg message.Message) error {
	log.Printf("RECV: %v\n", msg)
	return nil
}

func conn(srv *server.Server, conn server.PreConn, msg *message.ConnectMessage) error {
	log.Printf("CONNECT: %v\n", msg)
	conn.SetReceiveHandler(recv)
	return nil
}

func main() {
	srv := &server.Server{
		ConnectHandler: conn,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
