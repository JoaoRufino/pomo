package server

import (
	"fmt"
	"sync"

	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/server/rest"
	"github.com/joaorufino/pomo/pkg/server/unix"
)

// NewServer initializes the servers based on configuration
func NewServer(conf *conf.Config) (core.Server, error) {
	compositeServer := &CompositeServer{}

	if conf.Server.Type == "unix" || conf.Server.Type == "all" {
		unixServer := &unix.UnixServer{}
		server, err := unixServer.Init(conf)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Unix server: %w", err)
		}
		compositeServer.AddServer(server)
	}

	if conf.Server.Type == "rest" || conf.Server.Type == "all" {
		restServer := &rest.RestServer{}
		server, err := restServer.Init(conf)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize REST server: %w", err)
		}
		compositeServer.AddServer(server)
	}

	return compositeServer, nil
}

// CompositeServer holds multiple servers
type CompositeServer struct {
	servers []core.Server
	wg      sync.WaitGroup
}

// AddServer adds a new server to the composite server
func (cs *CompositeServer) AddServer(server core.Server) {
	cs.servers = append(cs.servers, server)
}

// Start starts all servers
func (cs *CompositeServer) Start() {
	cs.wg.Add(len(cs.servers))
	for _, server := range cs.servers {
		go func(s core.Server) {
			defer cs.wg.Done()
			s.Start()
		}(server)
	}
}

// Stop stops all servers
func (cs *CompositeServer) Stop() {
	for _, server := range cs.servers {
		server.Stop()
	}
	cs.wg.Wait()
}
