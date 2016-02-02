package main

import (
	"log"
	"strconv"

	"github.com/koron/go-mqtt/packet"
	"github.com/koron/go-mqtt/server"
)

type sadapter struct {
	nextID int
}

var _ server.Adapter = (*sadapter)(nil)

func (sa *sadapter) Connect(srv *server.Server, c server.Client, p *packet.Connect) (server.ClientAdapter, error) {
	log.Printf("CONNECT: %#v\n", p)
	id := sa.nextID
	sa.nextID++
	return &cadapter{
		srv: srv,
		cl:  c,
		id:  strconv.Itoa(id),
	}, nil
}

func (sa *sadapter) Disconnect(srv *server.Server, ca server.ClientAdapter, err error) {
	// nothing to do.
}

type cadapter struct {
	srv *server.Server
	cl  server.Client
	id  string
}

var (
	_ server.ClientAdapter = (*cadapter)(nil)
	_ server.PacketFilter  = (*cadapter)(nil)
)

func (ca *cadapter) ID() string {
	return ca.id
}

func (ca *cadapter) IsSessionPresent() bool {
	return false
}

func (ca *cadapter) OnDisconnect() error {
	log.Printf("DISCONNECTED: id=%s\n", ca.id)
	return nil
}

func (ca *cadapter) OnPing() (bool, error) {
	return true, nil
}

func (ca *cadapter) OnSubscribe(topics []server.Topic) ([]server.QoS, error) {
	qs := make([]server.QoS, len(topics))
	for i, t := range topics {
		log.Printf("SUBSCRIBE: %q qos=%s\n", t.Filter, t.QoS)
		qs[i] = server.AtMostOnce
	}
	return qs, nil
}

func (ca *cadapter) OnUnsubscribe(filters []string) error {
	// nothing to do.
	return nil
}

func (ca *cadapter) OnPublish(m *server.Message) error {
	log.Printf("PUBLISH: topic=%s body=%s (id=%s)\n", m.Topic, string(m.Body), ca.id)
	return nil
}

func (ca *cadapter) PreProcess(p packet.Packet) error {
	// nothing to do.
	return nil
}

func (ca *cadapter) PreSend(p packet.Packet, d []byte) ([]byte, error) {
	// nothing to do.
	return d, nil
}

func (ca *cadapter) PostSend(p packet.Packet, d []byte) {
	switch p.(type) {
	case *packet.SubACK:
		ca.cl.Publish(server.AtMostOnce, false, "a/b/c", []byte("Hello MQTT client"))
	}
}

func main() {
	srv := &server.Server{
		Adapter: new(sadapter),
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
