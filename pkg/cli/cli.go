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
	SetClient(*core.Client)
	SetConfig(*conf.Config)
}
