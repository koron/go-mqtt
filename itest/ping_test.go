package itest

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/koron/go-mqtt/client"
)

func TestPing(t *testing.T) {
	t.Parallel()
	srv := NewServer(t, nil, nil).Start()

	c := srv.Connect(t, client.Param{
		Options: &client.Options{
			CleanSession: true,
			KeepAlive:    60,
			Logger:       log.New(ioutil.Discard, "MQTT-C", log.LstdFlags),
		},
	})

	err := c.C.Ping()
	if err != nil {
		t.Fatalf("Ping() failed: %s", err)
	}
	c.Disconnect(t, false)

	srv.Stop()
}
