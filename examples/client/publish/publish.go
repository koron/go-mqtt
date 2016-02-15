package main

import (
	"log"
	"os"

	"github.com/koron/go-mqtt/client"
)

func main() {
	c, err := client.Connect(client.Param{
		ID: "example_publish",
		Options: &client.Options{
			CleanSession: true,
			KeepAlive:    60,
			Logger:       log.New(os.Stderr, "MQTT-C", log.LstdFlags),
		},
	})
	if err != nil {
		log.Fatalf("Connect() failed: %v", err)
	}
	err = c.Publish(client.AtMostOnce, false,
		"users/123/objects/789", []byte("Hello MQTT"))
	if err != nil {
		log.Printf("Publish() failed: %v", err)
	}
	err = c.Disconnect(false)
	if err != nil {
		log.Printf("Disconnect() failed: %v", err)
	}
	log.Print("DONE")
}
