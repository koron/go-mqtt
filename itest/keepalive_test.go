package itest

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/koron/go-mqtt/client"
	"github.com/koron/go-mqtt/server"
)

type syncErr struct {
	mu  sync.Mutex
	err error
}

func (se *syncErr) get() error {
	se.mu.Lock()
	defer se.mu.Unlock()
	return se.err
}

func (se *syncErr) set(err error) {
	se.mu.Lock()
	se.err = err
	se.mu.Unlock()
}

func TestKeepAlive(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	t.Parallel()

	var (
		keepAlive       = 2
		monitorInterval = 5
		pingExpect      = int(float64(monitorInterval) / (float64(keepAlive) - 0.5))
	)

	var pingMu sync.Mutex
	var pingCnt int
	srv := NewServer(t, &Adapter{
		onPing: func() (bool, error) {
			pingMu.Lock()
			pingCnt++
			pingMu.Unlock()
			return true, nil
		},
	}, nil).Start()

	ctx, cancel := context.WithCancel(context.Background())
	var cerr syncErr
	_ = srv.Connect(t, client.Param{
		Options: &client.Options{
			KeepAlive: uint16(keepAlive),
		},
		OnDisconnect: func(reason error, param client.Param) {
			cerr.set(reason)
			cancel()
		},
	})

	ti := time.NewTimer(time.Duration(monitorInterval) * time.Second)
wait:
	for {
		select {
		case <-ctx.Done():
			if !ti.Stop() {
				<-ti.C
			}
			break wait
		case <-ti.C:
			cancel()
			break wait
		}
	}

	if err := cerr.get(); err != nil {
		t.Fatalf("disconnected unexpectedly: %s", err)
		srv.Stop()
	}

	srv.Stop()

	if err := cerr.get(); err != nil {
		t.Fatal("client aliving unexpectedly")
	}

	pingMu.Lock()
	defer pingMu.Unlock()
	if pingCnt != pingExpect {
		t.Fatalf("ping count failed: expect=%d actual=%d", pingExpect, pingCnt)
	}
}

func TestKeepAlive_WithoutMonitor(t *testing.T) {
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
