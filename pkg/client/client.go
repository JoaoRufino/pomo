package client

import (
	"errors"

	"github.com/joaorufino/pomo/pkg/client/graphql"
	"github.com/joaorufino/pomo/pkg/client/rest"
	"github.com/joaorufino/pomo/pkg/client/unix"
	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
)

func NewClient(conf *conf.Config) (core.Client, error) {
	switch conf.Server.Type {
	case "unix":
		unixClient := &unix.UnixClient{}
		return unixClient.Init(conf)
	case "rest":
		restClient := &rest.RestClient{}
		return restClient.Init(conf)
	case "graphql":
		graphqlClient := &graphql.GraphQLClient{}
		return graphqlClient.Init(conf)
	}

	return nil, errors.New("error creating client")
}
