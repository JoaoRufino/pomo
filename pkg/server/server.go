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
}

type Runner interface {
	Status() *Status
}

func (s *Server) listen() {
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			break
		}
		buf := make([]byte, 512)
		// Ignore any content
		conn.Read(buf)
		raw, _ := json.Marshal(s.runner.Status())
		conn.Write(raw)
		conn.Close()
	}
}

func (s *Server) Start() {
	s.running = true
	go s.listen()
}

func (s *Server) Stop() {
	s.running = false
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

// Client makes requests to a listening
// pomo server to check the status of
// any currently running task session.
type Client struct {
	conn net.Conn
}

func (c Client) read(statusCh chan *Status) {
	buf := make([]byte, 512)
	n, _ := c.conn.Read(buf)
	status := &Status{}
	json.Unmarshal(buf[0:n], status)
	statusCh <- status
}

func (c Client) Status() (*Status, error) {
	statusCh := make(chan *Status)
	c.conn.Write([]byte("status"))
	go c.read(statusCh)
	return <-statusCh, nil
}

func (c Client) Close() error { return c.conn.Close() }

func NewClient(path string) (*Client, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}
