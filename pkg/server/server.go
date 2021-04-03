package server

import (
	"errors"
	"os"

	"github.com/joao.rufino/pomo/pkg/server/models"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

type Server interface {
	Start()
	Stop()
}

func NewServer(k *koanf.Koanf, runner models.Runner) (Server, error) {
	//check if socket file exists

	switch k.String("server.type") {
	case "unix":
		unix := &UnixServer{}
		return unix.Init(k.String("server.socket"), runner)
	}
	return nil, errors.New("unrecognized server type")

}

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:%s\n", err)
		os.Exit(1)
	}
}
