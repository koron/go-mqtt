package itest

import (
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/koron/go-mqtt/client"
	"github.com/koron/go-mqtt/server"
)

func newListener(tb testing.TB) net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Helper()
		tb.Fatalf("net.Listen failed: %s", err)
	}
	return l
}

// Server wraps *server.Server for test.
type Server struct {
	tb testing.TB
	l  net.Listener
	s  *server.Server
	wg sync.WaitGroup
	cn int
	mu sync.Mutex
}

// NewServer creates new server.Server for test.
func NewServer(tb testing.TB, a server.Adapter, opts *server.Options) *Server {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Helper()
		tb.Fatalf("net.Listen failed: %s", err)
	}
	return &Server{
		l: l,
		s: &server.Server{
			Addr:    "tcp://" + l.Addr().String(),
			Adapter: a,
			Options: opts,
		},
	}
}

// Start starts server in a new goroutine.
func (srv *Server) Start() *Server {
	srv.wg.Add(1)
	go func() {
		err := srv.s.Serve(srv.l)
		if err != nil {
			srv.tb.Helper()
			srv.tb.Errorf("server.Serve failed: %s", err)
		}
		srv.wg.Done()
	}()
	return srv
}

// Stop stops server and wait to terminate.
func (srv *Server) Stop() {
	srv.s.Close()
	srv.wg.Wait()
}

// Connect connects a client to test server.
func (srv *Server) Connect(tb testing.TB, p client.Param) *Client {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	c0 := &Client{}
	if p.Addr == "" {
		p.Addr = srv.s.Addr
	}
	if p.ID == "" {
		p.ID = fmt.Sprintf("%s-%d", tb.Name(), srv.cn)
		srv.cn++
		if n := len(p.ID); n > 20 {
			p.ID = p.ID[n-20:]
		}
	}
	c0.ID = p.ID
	if p.OnDisconnect == nil {
		p.OnDisconnect = c0.OnDisconnect
	}

	c, err := client.Connect(p)
	if err != nil {
		tb.Helper()
		tb.Fatalf("client.Connect failed: %s", err)
	}
	c0.C = c
	return c0
}

// Client is MQTT client object for test.
type Client struct {
	C  client.Client
	ID string
	mu sync.Mutex

	disconnectReason error
}

// OnDisconnect records last reason to disconnectReason.
func (c *Client) OnDisconnect(reason error, param client.Param) {
	c.mu.Lock()
	c.disconnectReason = reason
	c.mu.Unlock()
}

// DisconnectReason gets disconnect reason value.
func (c *Client) DisconnectReason() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.disconnectReason
}

// Disconnect disconnects from the server.
func (c *Client) Disconnect(tb testing.TB, force bool) {
	err := c.C.Disconnect(force)
	if err != nil {
		tb.Helper()
		tb.Fatalf("Client.Disconnect() failed: id=%s: %s", c.ID, err)
	}
}
