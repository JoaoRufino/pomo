package server

import (
	"fmt"

	"github.com/joao.rufino/pomo/pkg/cli"
	"github.com/spf13/cobra"
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewServerVersionCommand(pomoCli cli.Cli) *cobra.Command {
	serverVersionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  `Show version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version : " + pomoCli.Version() + " - " + pomoCli.Executable())
		},
	}
	return serverVersionCmd
}
