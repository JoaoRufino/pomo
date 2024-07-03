package server

import (
	"github.com/joaorufino/pomo/pkg/cli"
	"github.com/spf13/cobra"
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewServerConfigCommand(pomoCli cli.Cli) *cobra.Command {
	serverConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "Display the current server configuration",
		Long:  `Display the current server configuration`,
		Run: func(cmd *cobra.Command, args []string) {
			pomoCli.Config().Print()
		},
	}
	return serverConfigCmd
}
