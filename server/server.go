package server

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koron/go-mqtt/internal/backoff"
	"github.com/koron/go-mqtt/packet"
)

var (
	// ErrInvalidCloudAdapter indicates Adapter#Connect() returns invalid CloudAdapter.
	ErrInvalidCloudAdapter = errors.New("invalid CloudAdapter")

	// ErrUnknownProtocol indicates connect adddress includes unknown protocol.
	ErrUnknownProtocol = errors.New("unknown protocol")

	// ErrAlreadyServerd indicates the server is served already.
	ErrAlreadyServerd = errors.New("server served already")

	// ErrNotServing indicates the server is not under serving.
	ErrNotServing = errors.New("server is not serving")
)

const (
	none int32 = iota
	starting
	running
	closed
)

// Server is a instance of MQTT server.
type Server struct {
	Addr    string
	Adapter Adapter
	Options *Options

	st       int32
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
	var l net.Listener
	switch u.Scheme {
	case "tcp":
		l, err = net.Listen(u.Scheme, u.Host)
		if err != nil {
			return err
		}
	case "ssl", "tcps", "tls":
		l, err = tls.Listen("tcp", u.Host, srv.Options.TLSConfig)
		if err != nil {
			return err
		}
	default:
		return ErrUnknownProtocol
	}
	err = srv.Serve(l)
	if err != nil {
		if err == ErrAlreadyServerd {
			l.Close()
		}
		return err
	}
	return nil
}

// Serve accepts incoming connections on the Listener.
func (srv *Server) Serve(l net.Listener) error {
	if !atomic.CompareAndSwapInt32(&srv.st, none, starting) {
		return ErrAlreadyServerd
	}
	srv.logger = srv.options().Logger
	srv.quit = make(chan bool, 1)
	srv.listener = l
	srv.wg = sync.WaitGroup{}
	srv.cs = make(map[*client]bool)

	atomic.StoreInt32(&srv.st, running)
	srv.logServerStart()
	delay := backoff.Exp{Min: time.Millisecond * 5}
	for {
		conn, err := srv.listener.Accept()
		select {
		case <-srv.quit:
			go srv.terminateAllClients()
			return nil
		default:
		}
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				srv.logTemporaryError(nerr, &delay, nil)
				delay.Wait()
				continue
			}
			if atomic.CompareAndSwapInt32(&srv.st, running, closed) {
				srv.listener.Close()
			}
			go srv.terminateAllClients()
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
	if !atomic.CompareAndSwapInt32(&srv.st, running, closed) {
		return ErrNotServing
	}
	// terminate server.
	close(srv.quit)
	srv.listener.Close()
	srv.wg.Wait()
	return nil
}

func (srv *Server) terminateAllClients() {
	srv.cl.Lock()
	for c := range srv.cs {
		c.terminate()
	}
	srv.cl.Unlock()
}

func (srv *Server) logf(fmt string, a ...interface{}) {
	if srv.logger == nil {
		return
	}
	srv.logger.Printf(fmt, a...)
}

func (srv *Server) logServerStart() {
	srv.logf("MQTT server listen on: %s\n", srv.listener.Addr().String())
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
	srv.logf("failed to send packet;%#v to client;%s: %v", c.id(), p, err)
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
