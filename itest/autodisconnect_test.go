package itest

import (
	"testing"
	"time"

	"github.com/koron/go-mqtt/client"
)

func TestAutoDisconnect(t *testing.T) {
	t.Parallel()

	srv := NewServer(t, nil, nil).Start()

	c0 := srv.Connect(t, client.Param{
		Options: &client.Options{
			KeepAlive:            2,
			DisableAutoKeepAlive: true,
		},
	})

	// The connection is terminated once 1.5 times the KeepAlive interval -
	// i.e., 3 seconds - has elapsed. Here, an extra 100 milliseconds is
	// allowed to ensure reliability.  See [MQTT-3.1.2-24]
	time.Sleep(3*time.Second + 100*time.Millisecond)

	if c0.DisconnectReason() == nil {
		t.Error("not disconnected")
	}
	srv.Stop()
}
