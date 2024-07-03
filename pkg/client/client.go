package client

import (
	"errors"

	"github.com/joaorufino/pomo/pkg/client/rest"
	"github.com/joaorufino/pomo/pkg/client/unix"
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/spf13/viper"
)

func NewClient() (core.Client, error) {
	switch viper.GetString("server.type") {
	case "unix":
		unixClient := &unix.UnixClient{}
		return unixClient.Init(viper.GetViper())
	case "rest":
		restClient := &rest.RestClient{}
		return restClient.Init(viper.GetViper())
	}

	return nil, errors.New("error creating client")
}
