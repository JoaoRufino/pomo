package server

import (
	"encoding/json"
	"fmt"
	"net"
)

// Client makes requests to a listening
// pomo server to check the status of
// any currently running task session.
type Client struct {
	conn net.Conn
}

// read reads the status from the server
func (c *Client) read(statusCh chan *Status) {
	defer close(statusCh)

	buf := make([]byte, 512)
	n, err := c.conn.Read(buf)
	if err != nil {
		// Log the error
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	status := &Status{}
	if err := json.Unmarshal(buf[:n], status); err != nil {
		// Log the error
		fmt.Printf("Error unmarshalling status: %v\n", err)
		return
	}

	statusCh <- status
}

// Status requests the status from the server
func (c *Client) Status() (*Status, error) {
	statusCh := make(chan *Status)

	if _, err := c.conn.Write([]byte("status")); err != nil {
		return nil, fmt.Errorf("error writing to connection: %w", err)
	}

	go c.read(statusCh)
	return <-statusCh, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// NewClient initializes a new Client instance
func NewClient(path string) (*Client, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, fmt.Errorf("error dialing unix socket: %w", err)
	}
	return &Client{conn: conn}, nil
}
