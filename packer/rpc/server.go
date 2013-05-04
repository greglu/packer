package rpc

import (
	"errors"
	"github.com/mitchellh/packer/packer"
	"net"
	"net/rpc"
)

// A Server is a Golang RPC server that has helper methods for automatically
// setting up the endpoints for Packer interfaces.
type Server struct {
	listener net.Listener
	server *rpc.Server
}

// Creates and returns a new Server.
func NewServer() *Server {
	return &Server{
		server: rpc.NewServer(),
	}
}

func (s *Server) Address() string {
	if s.listener == nil {
		panic("Server not listening.")
	}

	return s.listener.Addr().String()
}

func (s *Server) RegisterUi(ui packer.Ui) {
	s.server.RegisterName("Ui", &UiServer{ui})
}

func (s *Server) Start() error {
	if s.listener != nil {
		return errors.New("Server already started.")
	}

	// Start the TCP listener and a goroutine responsible for cleaning up the
	// listener.
	s.listener = netListenerInRange(portRangeMin, portRangeMax)
	if s.listener == nil {
		return errors.New("Could not open a port ot listen on.")
	}

	// Start accepting connections
	go func(l net.Listener) {
		for {
			conn, err := l.Accept()
			if err != nil {
				break
			}

			go s.server.ServeConn(conn)
		}
	}(s.listener)

	return nil
}

func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}