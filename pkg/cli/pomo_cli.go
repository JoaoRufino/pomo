package cli

import (
	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
	"go.uber.org/zap"
)

type PomoCli struct {
	logger *zap.SugaredLogger
	config *conf.Config
	server *core.Server
	client core.Client
}

func NewPomoCli(configFile string) (*PomoCli, error) {
	config, err := conf.LoadConfig(configFile)
	if err != nil {
		return nil, err
	}

	conf.InitLogger()
	logger := zap.S().With("package", "cli")

	return &PomoCli{
		logger: logger,
		config: config,
	}, nil
}

func (pomoCli *PomoCli) Executable() string {
	return pomoCli.config.Server.Name
}

func (pomoCli *PomoCli) Version() string {
	return pomoCli.config.Server.Version
}

func (pomoCli *PomoCli) Config() *conf.Config {
	return pomoCli.config
}

func (pomoCli *PomoCli) Logger() *zap.SugaredLogger {
	return pomoCli.logger
}

func (pomoCli *PomoCli) Server() *core.Server {
	return pomoCli.server
}

func (pomoCli *PomoCli) Client() core.Client {
	return pomoCli.client
}

func (pomoCli *PomoCli) SetServer(server *core.Server) {
	pomoCli.server = server
}

func (pomoCli *PomoCli) SetConfig(config *conf.Config) {
	pomoCli.config = config
}

func (pomoCli *PomoCli) SetClient(client *core.Client) {
	if client != nil {
		pomoCli.client = *client
	} else {
		pomoCli.client = nil
	}
}

func (pomoCli *PomoCli) Close() {
	if pomoCli.client != nil {
		pomoCli.client.Close()
	}
}
