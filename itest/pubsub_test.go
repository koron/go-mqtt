package itest

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/koron/go-mqtt/client"
)

func TestPubSub(t *testing.T) {
	t.Parallel()
	srv := NewServer(t, &Adapter{}, nil).Start()

	c0 := srv.Connect(t, client.Param{
		Options: &client.Options{
			CleanSession: true,
			KeepAlive:    60,
			Logger:       log.New(ioutil.Discard, "MQTT-C0", log.LstdFlags),
		},
	})
	err := c0.C.Subscribe([]client.Topic{
		{Filter: "#", QoS: client.AtMostOnce},
	})
	if err != nil {
		log.Fatalf("c0.Subscribe() failed: %s", err)
	}
	ch := make(chan error)
	go func() {
		defer close(ch)
		m, err := c0.C.Read(true)
		if err != nil {
			ch <- fmt.Errorf("c0.C.Read() failed: %w", err)
			return
		}
		if !reflect.DeepEqual(m, &client.Message{
			Topic: "users/123/objects/789",
			Body:  []byte("Hello MQTT"),
		}) {
			ch <- fmt.Errorf("unexpected message: %+v", m)
			return
		}
	}()

	c1 := srv.Connect(t, client.Param{
		Options: &client.Options{
			CleanSession: true,
			KeepAlive:    60,
			Logger:       log.New(os.Stderr, "MQTT-C1", log.LstdFlags),
		},
	})
	err = c1.C.Publish(client.AtMostOnce, false, "users/123/objects/789", []byte("Hello MQTT"))
	if err != nil {
		t.Fatalf("c1.C.Publish() failed: %s", err)
	}

	// wait the goroutine ends
	err = <-ch
	if err != nil {
		t.Fatalf("read failure: %v", err)
	}

	c1.Disconnect(t, false)
	c0.Disconnect(t, false)

	srv.Stop()
}
