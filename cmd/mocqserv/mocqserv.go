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

func sent(conn server.Conn, msg message.Message) error {
	switch msg.(type) {
	case *message.SubackMessage:
		return conn.Publish("a/b/c", []byte("Hello MQTT"), message.QosAtMostOnce)
	}
	return nil
}

func sub(conn server.Conn, topic string, qos byte) (byte, error) {
	log.Printf("SUBSCRIBE: %q qos=%d\n", topic, qos)
	return qos, nil
}

func pub(conn server.Conn, msg *message.PublishMessage) error {
	log.Printf("PUBLISH: topic=%s payload=%s (id=%d)\n", string(msg.Topic()), string(msg.Payload()), conn.ID())
	return nil
}

func conn(srv *server.Server, conn server.PreConn, msg *message.ConnectMessage) error {
	log.Printf("CONNECT: %v\n", msg)
	conn.SetReceiveHandler(recv)
	conn.SetSentHandler(sent)
	conn.SetSubscribeHandler(sub)
	conn.SetPublishedHandler(pub)
	return nil
}

func disconn(srv *server.Server, conn server.DisConn, msg *message.DisconnectMessage) error {
	log.Printf("DISCONNECTED: %v\n", msg)
	return nil
}

func main() {
	srv := &server.Server{
		ConnectHandler:      conn,
		DisconnectedHandler: disconn,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
