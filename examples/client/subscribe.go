package main

import (
	"log"
	"os"
	"sync"

	"github.com/koron/go-mqtt/client"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	c1, err := client.Connect(client.Param{
		ID: "example_subscribe_1",
		Options: &client.Options{
			CleanSession: true,
			KeepAlive:    60,
			Logger:       log.New(os.Stderr, "MQTT-C1", log.LstdFlags),
		},
	}, nil)
	if err != nil {
		log.Fatalf("c1.Connect() failed: %v", err)
	}
	err = c1.Subscribe([]client.Topic{
		{Filter: "#", QoS: client.AtMostOnce},
	})
	if err != nil {
		log.Fatalf("c1.Subscribe() failed: %v", err)
	}

	go func() {
		defer wg.Done()
		m, err := c1.Read(true)
		if err != nil {
			log.Printf("c1.Read() failed: %v", err)
			return
		}
		if m == nil {
			log.Print("c1.Read() unexpected end")
			return
		}
		log.Printf("c1.Read() topic=%s body=%s", m.Topic, string(m.Body))
	}()

	c2, err := client.Connect(client.Param{
		ID: "example_subscribe_2",
		Options: &client.Options{
			CleanSession: true,
			KeepAlive:    60,
			Logger:       log.New(os.Stderr, "MQTT-C2", log.LstdFlags),
		},
	}, nil)
	if err != nil {
		log.Fatalf("c2.Connect() failed: %v", err)
	}
	err = c2.Publish(client.AtMostOnce, false,
		"users/123/objects/789", []byte("Hello MQTT"))
	if err != nil {
		log.Fatalf("c2.Publish() failed: %v", err)
	}

	wg.Wait()

	err = c2.Disconnect(false)
	if err != nil {
		log.Printf("c2.Disconnect() failed: %v", err)
	}
	err = c1.Disconnect(false)
	if err != nil {
		log.Printf("c1.Disconnect() failed: %v", err)
	}
	log.Print("DONE")
}
