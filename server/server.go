package server

import (
	"bufio"
	"errors"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/koron/go-debug"
	"github.com/surgemq/message"
)

var (
	// ErrServerClosed indicates a server have been closed already.
	ErrServerClosed = errors.New("server have been closed")

	// ErrTooManyConnections indicates too many connections for a server.
	ErrTooManyConnections = errors.New("too many connections")
)

// Server is a instance of MQTT server.
type Server struct {
	Addr    string
	Handler Handler

	quit     chan bool
	listener net.Listener

	connLock sync.Mutex
	conns    map[connID]*conn
	connID   connID // next connID to issue
}

func (srv *Server) addr() string {
	if srv.Addr == "" {
		return "tcp://127.0.0.1:1883"
	}
	return srv.Addr
}

// ListenAndServe listens on the TCP network address
func (srv *Server) ListenAndServe() error {
	u, err := url.Parse(srv.addr())
	if err != nil {
		return err
	}
	l, err := net.Listen(u.Scheme, u.Host)
	if err != nil {
		return err
	}
	srv.conns = make(map[connID]*conn)
	srv.connID = 0
	return srv.Serve(l)
}

// Serve accepts incomming connections on the Listener.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	debug.Printf("mqtt: start to listen on %s\n", l.Addr().String())
	srv.quit = make(chan bool, 1)
	srv.listener = l
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, err := l.Accept()
		select {
		case <-srv.quit:
			return nil
		default:
		}
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.logf("mqtt: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		c := newConn(srv, rw)
		go c.serve()
	}
}

// Close terminates the server by shutting down all the client connections and
// closing.
func (srv *Server) Close() error {
	close(srv.quit)
	srv.listener.Close()
	// close all connections.
	srv.connLock.Lock()
	defer srv.connLock.Unlock()
	for _, c := range srv.conns {
		c.Close()
	}
	srv.conns = nil
	return nil
}

func (srv *Server) logf(format string, args ...interface{}) {
	debug.Printf(format, args...)
}

func (srv *Server) newConn(rwc net.Conn) (*conn, error) {
	c := &conn{
		server: srv,
		rwc:    rwc,
		reader: bufio.NewReader(rwc),
		writer: rwc,
	}
	return c, nil
}

func (srv *Server) authenticate(c *conn, m *message.ConnectMessage) error {
	// TODO: authenticate ConnectMessage.
	return nil
}

func (srv *Server) register(c *conn) error {
	srv.connLock.Lock()
	defer srv.connLock.Unlock()
	if srv.conns == nil {
		return ErrServerClosed
	}
	start := srv.connID
	for {
		if _, ok := srv.conns[srv.connID]; !ok && srv.connID != 0 {
			break
		}
		srv.connID++
		if start == srv.connID {
			return ErrTooManyConnections
		}
	}
	c.id = srv.connID
	srv.conns[srv.connID] = c
	srv.connID++
	return nil
}

func (srv *Server) unregister(c *conn) {
	srv.connLock.Lock()
	defer srv.connLock.Unlock()
	if srv.conns == nil {
		return
	}
	v, ok := srv.conns[c.id]
	if !ok || v != c {
		return
	}
	delete(srv.conns, c.id)
}
