package cli

import (
	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
	"go.uber.org/zap"
)

// Cli represents the pomo command line.
// It shall be implemented by root and by the mockCLi test package
type Cli interface {
	Version() string
	Executable() string
	Config() *conf.Config
	Logger() *zap.SugaredLogger
	Server() *core.Server
	Client() core.Client
	SetServer(*core.Server)
	SetClient(core.Client)
}

// PomoCli is an instance the docker command line client.
// Instances of the client can be returned from NewPomoCli.
type PomoCli struct {
	logger *zap.SugaredLogger
	config *conf.Config
	server *core.Server
	client core.Client
}

// NewPomoCli creates a new instance of PomoCli.
func NewPomoCli(configFile string) (*PomoCli, error) {
	pomoCli := &PomoCli{}

	// Load configuration
	var err error
	if configFile != "" {
		pomoCli.config, err = conf.LoadConfig(configFile)
		if err != nil {
			return nil, err
		}
	} else {
		pomoCli.config = conf.LoadDefaultConfig()
	}

	// Initialize the logger
	conf.InitLogger()

	pomoCli.logger = zap.S().With("package", "cli")

	return pomoCli, nil
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

func (pomoCli *PomoCli) SetClient(client core.Client) {
	pomoCli.client = client
}
