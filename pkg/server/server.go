package server

import (
	"errors"
	"log"
	"os"

	"github.com/joao.rufino/pomo/pkg/core"
	"github.com/joao.rufino/pomo/pkg/core/models"
	"github.com/joao.rufino/pomo/pkg/server/rest"
	"github.com/joao.rufino/pomo/pkg/server/unix"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

func NewServer(k *koanf.Koanf, runner models.Runner) (core.Server, error) {
	//check if socket file exists

	switch k.String("server.type") {
	case "unix":
		server := &unix.UnixServer{}
		return server.Init(k, runner)
	case "rest":
		s, err := rest.New(k)
		if err != nil {
			log.Fatalf("Could not create server", "error", err)
		}
		return s, nil
	}
	return nil, errors.New("unrecognized server type")
}

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:%s\n", err)
		os.Exit(1)
	}
}
