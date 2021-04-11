package client

import (
	"errors"

	"github.com/joao.rufino/pomo/pkg/client/rest"
	"github.com/joao.rufino/pomo/pkg/client/unix"
	"github.com/joao.rufino/pomo/pkg/core"
	"github.com/knadh/koanf"
)

func NewClient(k *koanf.Koanf) (core.Client, error) {
	switch k.String("server.type") {
	case "unix":
		unixClient := &unix.UnixClient{}
		return unixClient.Init(k)
	case "rest":
		restClient := &rest.RestClient{}
		return restClient.Init(k)
	}

	return nil, errors.New("error creating client")
}
