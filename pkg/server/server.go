package server

import (
	"errors"
	"os"

	"github.com/joao.rufino/pomo/pkg/core"
	"github.com/joao.rufino/pomo/pkg/core/models"
	"github.com/joao.rufino/pomo/pkg/server/unix"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

func NewServer(k *koanf.Koanf, runner models.Runner) (core.Server, error) {
	//check if socket file exists

	switch k.String("server.type") {
	case "unix":
		unix := &unix.UnixServer{}
		return unix.Init(k, runner)
	}
	return nil, errors.New("unrecognized server type")

}

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:%s\n", err)
		os.Exit(1)
	}
}
