package server

import (
	"errors"

	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/server/rest"
	"github.com/joaorufino/pomo/pkg/server/unix"
)

func NewServer(conf *conf.Config) (core.Server, error) {
	switch conf.Server.Type {
	case "unix":
		unixServer := &unix.UnixServer{}
		return unixServer.Init(conf)
	case "rest":
		restServer := &rest.RestServer{}
		return restServer.Init(conf)
	}

	return nil, errors.New("error creating server")
}
