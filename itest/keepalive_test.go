package itest

import (
	"sync"
	"testing"
	"time"

	"github.com/koron/go-mqtt/client"
	"github.com/koron/go-mqtt/server"
)

func TestKeepAlive(t *testing.T) {
	var wg sync.WaitGroup
	srv := server.Server{
		Options: &server.Options{
			DisableMonitor: true,
		},
	}
	l := newListener(t)

	wg.Add(1)
	go func() {
		err := srv.Serve(l)
		if err != nil {
			t.Fatal("Serve is failed: ", err)
		}
		wg.Done()
	}()

	time.Sleep(time.Millisecond * 100)

	var disconnReason error
	var mu sync.Mutex
	_, err := client.Connect(client.Param{
		Addr: "tcp://" + l.Addr().String(),
		OnDisconnect: func(reason error, param client.Param) {
			mu.Lock()
			disconnReason = reason
			mu.Unlock()
		},
		ID: t.Name(),
		Options: &client.Options{
			KeepAlive:            2,
			DisableAutoKeepAlive: true,
		},
	})
	if err != nil {
		t.Fatal("client.Connect failed: ", err)
	}

	time.Sleep(time.Second * 3)

	mu.Lock()
	if disconnReason != nil {
		t.Error("disconnected unexpectedly: ", disconnReason)
	}
	mu.Unlock()
	srv.Close()
	wg.Wait()

	time.Sleep(time.Millisecond * 100)
	mu.Lock()
	if disconnReason == nil {
		t.Error("client aliving unexpectedly")
	}
	mu.Unlock()
}
