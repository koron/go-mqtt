package server

import (
	"bufio"
	"io"
	"net"
	"net/url"
	"time"

	"github.com/koron/go-debug"
	"github.com/surgemq/message"
)

// Server is a instance of MQTT server.
type Server struct {
	Addr    string
	Handler Handler

	quit     chan bool
	listener net.Listener
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
		c, err := srv.newConn(rw)
		if err != nil {
			debug.Printf("newConn failed: %v\n", err)
			continue
		}
		go c.serve()
	}
}

// Close terminates the server by shutting down all the client connections and
// closing.
func (srv *Server) Close() error {
	close(srv.quit)
	srv.listener.Close()
	// TODO: close each connections.
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

type conn struct {
	server *Server
	rwc    net.Conn
	reader *bufio.Reader
	writer io.Writer
}

func (c *conn) Close() error {
	c.rwc.Close()
	// FIXME: notify closing to server.
	return nil
}

func (c *conn) establishConnection() error {
	req, err := readConnectMessage(c.reader)
	if err != nil {
		writeConnackErrorMessage(c.writer, err)
		return err
	}
	err = c.server.authenticate(c, req)
	if err != nil {
		writeConnackErrorMessage(c.writer, err)
		return err
	}
	// send connack message.
	resp := message.NewConnackMessage()
	resp.SetSessionPresent(true)
	resp.SetReturnCode(message.ConnectionAccepted)
	err = writeMessage(c.writer, resp)
	if err != nil {
		return err
	}
	return nil
}

func (c *conn) serve() {
	defer c.Close()
	err := c.establishConnection()
	if err != nil {
		debug.Printf("establishConnection failed: %v\n", err)
		return
	}
	go c.recvMain()
	c.sendMain()
}

func (c *conn) recvMain() {
	for {
		// TODO:
	}
}

func (c *conn) sendMain() {
	for {
		// TODO:
	}
}
