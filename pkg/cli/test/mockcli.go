package test

import (
	"github.com/joao.rufino/pomo/pkg/client/test"
	"github.com/joao.rufino/pomo/pkg/conf"
	"github.com/joao.rufino/pomo/pkg/core"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

type MockCli struct {
	version    string
	executable string
	config     *koanf.Koanf
	logger     *zap.SugaredLogger
	client     core.Client
	server     *core.Server
}

func NewMockCli() *MockCli {
	cli := &MockCli{}

	// Global koanf configuration with "." for delimeter
	cli.config = koanf.New(".")
	_ = ConfFromDefaults(cli.config)

	client := test.NewMockClient(cli.config, test.MockClientOptions{})
	cli.client = client

	conf.InitLogger(cli.Config())
	cli.logger = zap.S().With("package", "test")
	return cli
}

func (cli *MockCli) Version() string {
	return cli.version
}

func (cli *MockCli) SetVersion(version string) {
	cli.version = version
}

func (cli *MockCli) Executable() string {
	return cli.executable
}

func (cli *MockCli) SetExecutable(executable string) {
	cli.executable = executable
}

func (cli *MockCli) Config() *koanf.Koanf {
	return cli.config
}

func (cli *MockCli) Logger() *zap.SugaredLogger {
	return cli.logger

}

// SetLogger sets the "fake" logger
func (cli *MockCli) SetLogger(logger *zap.SugaredLogger) {
	cli.logger = logger
}

// SetConfigFile sets the "fake" config file
func (cli *MockCli) SetConfigFile(configfile *koanf.Koanf) {
	cli.config = configfile
}

func (cli *MockCli) Client() core.Client {
	return cli.client
}

func (cli *MockCli) Server() *core.Server {
	return cli.server
}
func (cli *MockCli) SetServer(server *core.Server) {
	cli.server = server
}
func (cli *MockCli) SetClient(client *core.Client) {
	cli.client = *client
}
