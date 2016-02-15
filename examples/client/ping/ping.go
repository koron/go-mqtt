package main

import (
	"log"
	"os"

	"github.com/koron/go-mqtt/client"
)

func main() {
	c, err := client.Connect(client.Param{
		ID: "example_ping",
		Options: &client.Options{
			CleanSession: true,
			KeepAlive:    60,
			Logger:       log.New(os.Stderr, "MQTT-C", log.LstdFlags),
		},
	})
	if err != nil {
		log.Fatalf("Connect() failed: %v", err)
	}
	err = c.Ping()
	if err != nil {
		log.Printf("Ping() failed: %v", err)
	}
	err = c.Disconnect(false)
	if err != nil {
		log.Printf("Disconnect() failed: %v", err)
	}
	log.Print("DONE")
}
