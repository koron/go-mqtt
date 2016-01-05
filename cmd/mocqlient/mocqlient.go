package main

import (
	"log"

	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

func main() {
	c := &service.Client{}

	msgConn := message.NewConnectMessage()
	msgConn.SetVersion(3)
	msgConn.SetCleanSession(true)
	msgConn.SetClientId([]byte("01234567"))
	msgConn.SetKeepAlive(300)

	if err := c.Connect("tcp://127.0.0.1:1883", msgConn); err != nil {
		log.Fatal(err)
	}
	log.Println("connected")

	msgSub := message.NewSubscribeMessage()
	msgSub.AddTopic([]byte("#"), 0)
	onComp := func(msg, ack message.Message, err error) error {
		log.Printf("on complete subscribe: msg=%#v ack=%#v err=%v\n", msg, ack, err)
		return err
	}
	onPub := func(msg *message.PublishMessage) error {
		log.Printf("on publish: %#v\n", msg)
		return nil
	}
	if err := c.Subscribe(msgSub, onComp, onPub); err != nil {
		log.Fatal(err)
	}
	log.Println("subscribed")

	//c.Disconnect()
	//log.Println("disconnected")

	done := make(chan bool)
	<-done
}
