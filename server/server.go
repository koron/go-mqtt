package server

import (
	"bufio"
	"errors"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/koron/go-debug"
	"github.com/koron/go-mqtt/internal/backoff"
	"github.com/koron/go-mqtt/packet"
)

var (
	// ErrServerClosed indicates a server have been closed already.
	ErrServerClosed = errors.New("server have been closed")

	// ErrTooManyConnections indicates too many connections for a server.
	ErrTooManyConnections = errors.New("too many connections")

	// ErrInvalidCloudAdapter indicates Adapter#Connect() returns invalid CloudAdapter.
	ErrInvalidCloudAdapter = errors.New("invalid CloudAdapter")
)

// Server is a instance of MQTT server.
type Server struct {
	Addr    string
	Adapter Adapter

	quit     chan bool
	listener net.Listener
	wg       sync.WaitGroup // for Client#serve()

	connLock sync.Mutex
	conns    map[ConnID]*conn
	connID   ConnID         // next ConnID to issue
	connWait sync.WaitGroup // wait goroutines of conn.serve()
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
	srv.conns = make(map[ConnID]*conn)
	srv.connID = 0
	srv.connWait = sync.WaitGroup{}
	return srv.Serve(l)
}

// Serve accepts incomming connections on the Listener.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	debug.Printf("mqtt: start to listen on %s\n", l.Addr().String())
	srv.quit = make(chan bool, 1)
	srv.listener = l
	srv.wg = sync.WaitGroup{}
	delay := backoff.Exp{Min: time.Millisecond * 5}
	for {
		conn, err := l.Accept()
		select {
		case <-srv.quit:
			return nil
		default:
		}
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				srv.logTemporaryError(nerr, &delay, nil)
				delay.Wait()
				continue
			}
			return err
		}
		delay.Reset()

		// start client goroutine.
		c := newClient(srv, conn)
		srv.wg.Add(1)
		go func() {
			c.serve()
			srv.wg.Done()
		}()
	}
}

// Close terminates the server by shutting down all the client connections and
// closing.
func (srv *Server) Close() error {
	close(srv.quit)
	srv.listener.Close()
	// close all connections.
	srv.connLock.Lock()
	for _, c := range srv.conns {
		c.Close()
	}
	srv.conns = nil
	srv.connLock.Unlock()
	// wait to terminate all clients.
	srv.wg.Wait()
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

func (srv *Server) logTemporaryError(err net.Error, d *backoff.Exp, c *client) {
	// TODO: log net temporary error.
}

func (srv *Server) logEstablishFailure(c *client, err error) {
	// TODO: log establish failure.
}

func (srv *Server) adapter() Adapter {
	if srv.Adapter == nil {
		return DefaultAdapter
	}
	return srv.Adapter
}

func (srv *Server) clientOnConnect(c *client, p *packet.Connect) (ClientAdapter, error) {
	ca, err := srv.adapter().Connect(srv, p)
	if err != nil {
		return nil, err
	}
	if ca == nil {
		return nil, ErrInvalidCloudAdapter
	}
	return ca, nil
}

func (srv *Server) clientOnStart(c *client) {
	// TODO: register living client.
}

func (srv *Server) clientOnStop(c *client, err error) {
	srv.adapter().Disconnect(srv, c.ca, err)
	// TODO: unregister client
}

func (srv *Server) clientOnDisconnect(c *client) {
	// TODO:
}
