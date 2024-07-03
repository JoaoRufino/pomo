package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/core/models"
	"github.com/joaorufino/pomo/pkg/server/rest"
)

// Session listens on a Unix domain socket
// for Pomo status requests
type Session struct {
	listener net.Listener
	runner   models.Runner
	running  bool
	mu       sync.Mutex
}

// listen handles incoming connections and responds with the current status
func (s *Session) listen() {
	for {
		s.mu.Lock()
		if !s.running {
			s.mu.Unlock()
			break
		}
		s.mu.Unlock()

		conn, err := s.listener.Accept()
		if err != nil {
			// Log the error and continue
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection processes a single connection
func (s *Session) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 512)
	if _, err := conn.Read(buf); err != nil {
		// Log the error
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	status := s.runner.Status()
	raw, err := json.Marshal(status)
	if err != nil {
		// Log the error
		fmt.Printf("Error marshalling status: %v\n", err)
		return
	}

	if _, err := conn.Write(raw); err != nil {
		// Log the error
		fmt.Printf("Error writing to connection: %v\n", err)
	}
}

// Start begins the server's listening process
func (s *Session) Start() {
	s.mu.Lock()
	s.running = true
	s.mu.Unlock()
	go s.listen()
}

// Stop halts the server's listening process
func (s *Session) Stop() {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
	s.listener.Close()
}

// NewServer creates a new server based on the configuration
func NewServer(config *conf.Config, runner models.Runner) (core.Server, error) {
	switch config.Server.Type {
	case "unix":
		listener, err := net.Listen("unix", config.Server.UnixSocket)
		if err != nil {
			return nil, err
		}
		session := &Session{listener: listener, runner: runner}
		session.Start()
		return session, nil

	case "rest":
		s, err := rest.New(config)
		if err != nil {
			log.Fatalf("Could not create server: %v", err)
		}
		return s, nil

	default:
		return nil, fmt.Errorf("unknown server type: %s", config.Server.Type)
	}
}
