package main

import (
	"log"

	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

func main() {
	c := &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetVersion(3)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte("01234567"))
	msg.SetKeepAlive(300)

	if err := c.Connect("tcp://127.0.0.1:1883", msg); err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	<-done
}
