package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	"github.com/koron/go-mqtt/client"
)

func onPublish(m *client.Message) {
	log.Printf("onPublish: {Topic:%q, Body:%q}", m.Topic, string(m.Body))
}

func loadCertPool(path string) (*x509.CertPool, error) {
	cp := x509.NewCertPool()
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cp.AppendCertsFromPEM(d)
	return cp, nil
}

func main() {
	cp, err := loadCertPool("../server.crt")
	if err != nil {
		log.Fatal("loadCertPool failed: ", err)
	}
	c, err := client.Connect(client.Param{
		Addr:      "tls://127.0.0.1:8883",
		ID:        "tlstest",
		OnPublish: onPublish,
		Options: &client.Options{
			Version:      4,
			CleanSession: true,
			KeepAlive:    60,
			TLSConfig: &tls.Config{
				RootCAs:            cp,
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		log.Fatal("client.Connect failed: ", err)
	}
	defer func() {
		err := c.Disconnect(false)
		log.Print("disconnected: ", err)
	}()
	log.Print("connected")
	err = c.Subscribe([]client.Topic{
		{Filter: "#", QoS: client.AtMostOnce},
	})
	if err != nil {
		log.Print("c.Subscribe failed: ", err)
	}
	log.Print("subscribed")
	time.Sleep(time.Second * 10)
}
