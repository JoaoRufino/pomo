package server

import (
	"fmt"

	"github.com/joao.rufino/pomo/pkg/conf"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewServerVersionCommand(cmd *cli.Command) *cobra.Command {
	serverVersionCmd := &cli.Command{
		Use:   "version",
		Short: "Show version",
		Long:  `Show version`,
		Run: func(cmd *cli.Command, args []string) {
			fmt.Println("Version : " + conf.K.String("server.name") + " - " + conf.K.String("server.version"))
		},
	}
	return serverVersionCmd
}
