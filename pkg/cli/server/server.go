package server

import (
	"github.com/joao.rufino/pomo/pkg/cli"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

// server command
//  pomo
//   ├── server
//   │   ├── config
//   │   ├── init
//   |   ├── status
//   │   └── version
///

// NewServerCommand returns a cobra command for `server` subcommands
func NewServerCommand(pomoCli cli.Cli) *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "operations regarding the server",
		Long:  "operations affecting the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	serverCmd.AddCommand(
		NewServerConfigCommand(pomoCli),
		NewServerStatusCommand(pomoCli),
		NewServerInitCommand(pomoCli),
		NewServerVersionCommand(pomoCli),
	)
	return serverCmd
}

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:%s", err)
	}
}
