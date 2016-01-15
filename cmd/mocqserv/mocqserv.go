package main

import (
	"log"

	"github.com/koron/go-mqtt/packet"
	"github.com/koron/go-mqtt/server"
)

func recv(conn server.Conn, p packet.Packet) error {
	log.Printf("RECV: %#v\n", p)
	return nil
}

func sent(conn server.Conn, p packet.Packet) error {
	switch p.(type) {
	case *packet.SubACK:
		return conn.Publish("a/b/c", []byte("Hello MQTT"), packet.QAtMostOnce)
	}
	return nil
}

func onSub(conn server.Conn, t packet.Topic) (packet.SubscribeResult, error) {
	log.Printf("SUBSCRIBE: %q qos=%d\n", t.Filter, t.RequestedQoS)
	return packet.SubscribeResult(t.RequestedQoS), nil
}

func onPub(conn server.Conn, p *packet.Publish) error {
	log.Printf("PUBLISH: topic=%s payload=%s (id=%d)\n", p.TopicName, string(p.Payload), conn.ID())
	return nil
}

func onConn(srv *server.Server, conn server.PreConn, p *packet.Connect) packet.ConnectReturnCode {
	log.Printf("CONNECT: %#v\n", p)
	conn.SetReceiveHandler(recv)
	conn.SetSentHandler(sent)
	conn.SetSubscribeHandler(onSub)
	conn.SetPublishedHandler(onPub)
	return packet.ConnectAccept
}

func onDisconn(srv *server.Server, conn server.DisConn, p *packet.Disconnect) error {
	log.Printf("DISCONNECTED: %#v\n", p)
	return nil
}

func main() {
	srv := &server.Server{
		ConnectHandler:      onConn,
		DisconnectedHandler: onDisconn,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
