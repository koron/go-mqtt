package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/koron/go-mqtt/packet"
	"github.com/koron/go-mqtt/server"
)

type sadapter struct {
	clients map[string]*server.NullClientAdapter
}

func (sa *sadapter) Connect(srv *server.Server, c server.Client, p *packet.Connect) (server.ClientAdapter, error) {
	ca := &server.NullClientAdapter{
		Client:   c,
		ClientID: p.ClientID,
	}
	sa.clients[ca.ID()] = ca
	return ca, nil
}

func (sa *sadapter) Disconnect(srv *server.Server, ca server.ClientAdapter, err error) {
	if ca == nil {
		return
	}
	delete(sa.clients, ca.ID())
}

func (sa *sadapter) SendToAll(topic string, body []byte) {
	for _, ca := range sa.clients {
		go ca.Client.Publish(server.AtMostOnce, false, topic, body)
	}
}

func loadTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	ca, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{ca},
	}, nil
}

func main() {
	tc, err := loadTLSConfig("../server.crt", "../server.key")
	if err != nil {
		log.Fatal("loadTLSConfig failed: ", err)
	}
	sa := &sadapter{
		clients: map[string]*server.NullClientAdapter{},
	}
	s := &server.Server{
		Addr:    "tls://127.0.0.1:8883",
		Adapter: sa,
		Options: &server.Options{
			TLSConfig: tc,
		},
	}
	go func() {
		count := 0
		for {
			time.Sleep(time.Second)
			m := fmt.Sprintf("message #%d", count)
			sa.SendToAll("count", []byte(m))
			count++
			log.Printf("sent: %s", m)

		}
		fmt.Println("HERE")
	}()
	s.ListenAndServe()
}
