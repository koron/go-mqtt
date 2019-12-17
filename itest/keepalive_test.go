package itest

import (
	"testing"
	"time"

	"github.com/koron/go-mqtt/client"
	"github.com/koron/go-mqtt/server"
)

func TestKeepAlive(t *testing.T) {
	t.Parallel()

	srv := NewServer(t, nil, &server.Options{
		DisableMonitor: true,
	}).Start()

	c0 := srv.Connect(t, client.Param{
		Options: &client.Options{
			KeepAlive:            2,
			DisableAutoKeepAlive: true,
		},
	})

	time.Sleep(time.Second * 3)

	if err := c0.DisconnectReason(); err != nil {
		t.Errorf("disconnected unexpectedly: %s", err)
	}

	srv.Stop()

	time.Sleep(time.Millisecond * 100)
	if err := c0.DisconnectReason(); err == nil {
		t.Error("client aliving unexpectedly")
	}
}
