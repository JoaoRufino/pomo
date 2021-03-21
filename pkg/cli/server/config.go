package server

import (
	"github.com/joao.rufino/pomo/pkg/conf"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewServerConfigCommand(cmd *cli.Command) *cobra.Command {
	serverConfigCmd := &cli.Command{
		Use:   "config",
		Short: "Display the current server configuration",
		Long:  `Display the current server configuration`,
		Run: func(cmd *cli.Command, args []string) {
			conf.K.Print()
		},
	}
	return serverConfigCmd
}
