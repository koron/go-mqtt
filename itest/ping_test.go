package itest

import (
	"log"
	"os"
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
			Logger:       log.New(os.Stderr, "MQTT-C", log.LstdFlags),
		},
	})

	err := c.C.Ping()
	if err != nil {
		t.Fatalf("Ping() failed: %s", err)
	}
	err = c.C.Disconnect(false)
	if err != nil {
		t.Fatalf("Disconnect() faield: %s", err)
	}

	srv.Stop()
}
