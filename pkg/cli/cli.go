package cli

import (
	"github.com/joao.rufino/pomo/pkg/client"
	"github.com/joao.rufino/pomo/pkg/conf"
	"github.com/joao.rufino/pomo/pkg/core"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

// Cli represents the pomo command line.
// It shall be implemented by root and by the mockCLi test package
type Cli interface {
	Version() string
	Executable() string
	Config() *koanf.Koanf
	Logger() *zap.SugaredLogger
	Server() *core.Server
	Client() core.Client
	SetServer(*core.Server)
}

// PomoCli is an instance the docker command line client.
// Instances of the client can be returned from NewPomoCli.
type PomoCli struct {
	logger *zap.SugaredLogger
	config *koanf.Koanf
	server *core.Server
	client core.Client
}

// constructor for abstract class
func NewPomoCli() (*PomoCli, error) {
	pomoCli := &PomoCli{}

	// Global koanf configuration with "." for delimeter
	pomoCli.config = koanf.New(".")

	// Load configuration
	err := conf.ConfFromDefaults(pomoCli.config)
	if err != nil {
		return nil, err
	}
	conf.InitLogger(pomoCli.Config())
	pomoCli.logger = zap.S().With("package", "cli")
	pomoCli.client, err = client.NewClient(pomoCli.Config())
	if err != nil {
		return nil, err
	}

	return pomoCli, nil
}

func (pomoCli *PomoCli) Executable() string {
	return pomoCli.config.String("server.name")
}

func (pomoCli *PomoCli) Version() string {
	return pomoCli.config.String("server.version")
}

func (pomoCli *PomoCli) Config() *koanf.Koanf {
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
