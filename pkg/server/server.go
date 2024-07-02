package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/joao.rufino/pomo/pkg/conf"
)

// Server listens on a Unix domain socket
// for Pomo status requests
type Server struct {
	listener net.Listener
	runner   Runner
	running  bool
	mu       sync.Mutex
}

type Runner interface {
	Status() *Status
}

// listen handles incoming connections and responds with the current status
func (s *Server) listen() {
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
func (s *Server) handleConnection(conn net.Conn) {
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
func (s *Server) Start() {
	s.mu.Lock()
	s.running = true
	s.mu.Unlock()
	go s.listen()
}

// Stop halts the server's listening process
func (s *Server) Stop() {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
	s.listener.Close()
}


func NewServer(runner Runner) (*Server, error) {
	//check if socket file exists

	socketPath := conf.K.String("server.socket")

	if _, err := os.Stat(socketPath); err == nil {
		_, err := net.Dial("unix", socketPath)
		//if error then sock file was saved after crash
		if err != nil {
			os.Remove(socketPath)
		} else {
			// another instance of pomo is running
			return nil, fmt.Errorf("socket %s is already in use", socketPath)
		}
	}
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}
	return &Server{listener: listener, runner: runner}, nil
}

