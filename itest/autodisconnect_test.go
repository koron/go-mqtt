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

	time.Sleep(time.Second * 3)

	if c0.DisconnectReason() == nil {
		t.Error("not disconnected")
	}
	srv.Stop()
}
