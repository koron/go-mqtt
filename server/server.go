package server

import (
	"errors"
	"log"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/koron/go-mqtt/internal/backoff"
	"github.com/koron/go-mqtt/packet"
)

var (
	// ErrInvalidCloudAdapter indicates Adapter#Connect() returns invalid CloudAdapter.
	ErrInvalidCloudAdapter = errors.New("invalid CloudAdapter")
)

// Server is a instance of MQTT server.
type Server struct {
	Addr    string
	Adapter Adapter
	Options *Options

	logger   *log.Logger
	quit     chan bool
	listener net.Listener
	wg       sync.WaitGroup // for client#serve()
	cl       sync.Mutex
	cs       map[*client]bool
}

func (srv *Server) addr() string {
	if srv.Addr == "" {
		return "tcp://127.0.0.1:1883"
	}
	return srv.Addr
}

func (srv *Server) options() *Options {
	if srv.Options == nil {
		return DefaultOptions
	}
	return srv.Options
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
	return srv.Serve(l)
}

// Serve accepts incoming connections on the Listener.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	srv.logger = srv.options().Logger
	srv.quit = make(chan bool, 1)
	srv.listener = l
	srv.wg = sync.WaitGroup{}
	srv.cs = make(map[*client]bool)
	srv.logServerStart(l)
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
	// terminate server.
	close(srv.quit)
	srv.listener.Close()
	// terminate all clients
	srv.cl.Lock()
	for c := range srv.cs {
		c.terminate()
	}
	srv.cl.Unlock()
	// wait to terminate all clients.
	srv.wg.Wait()
	return nil
}

func (srv *Server) logf(fmt string, a ...interface{}) {
	if srv.logger == nil {
		return
	}
	srv.logger.Printf(fmt, a...)
}

func (srv *Server) logServerStart(l net.Listener) {
	srv.logf("server starts to listen: %s\n", l.Addr().String())
}

func (srv *Server) logTemporaryError(err net.Error, d *backoff.Exp, c *client) {
	if c != nil {
		srv.logf("client;%s detect temporary error: %v", c.id(), err)
		return
	}
	srv.logf("server detect temporary error: %v", err)
}

// logAdapterError logs continuable AdapterError
func (srv *Server) logAdapterError(err AdapterError, p packet.Packet, c *client) {
	srv.logf("client:%s rejects packet:%#v but continue: %s",
		c.id(), p, err.Error())
}

func (srv *Server) logEstablishFailure(c *client, err error) {
	srv.logf("client fails to connect: %v", err)
}

func (srv *Server) logSendPacketError(c *client, p packet.Packet, err error) {
	srv.logf("failed to send packet;%#v to client;%s: %v", c.id, p, err)
}

func (srv *Server) adapter() Adapter {
	if srv.Adapter == nil {
		return DefaultAdapter
	}
	return srv.Adapter
}

func (srv *Server) clientOnConnect(c *client, p *packet.Connect) (ClientAdapter, error) {
	ca, err := srv.adapter().Connect(srv, c, p)
	if err != nil {
		return nil, err
	}
	if ca == nil {
		return nil, ErrInvalidCloudAdapter
	}
	return ca, nil
}

func (srv *Server) clientOnStart(c *client) {
	srv.cl.Lock()
	defer srv.cl.Unlock()
	srv.cs[c] = true
}

func (srv *Server) clientOnStop(c *client) {
	srv.cl.Lock()
	defer srv.cl.Unlock()
	delete(srv.cs, c)
}

func (srv *Server) clientOnDisconnect(c *client, err error) {
	srv.adapter().Disconnect(srv, c.ca, err)
}
