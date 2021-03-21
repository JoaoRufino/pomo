package server

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
)

// server command
//  pomo
//   ├── server
//   │   ├── config
//   │   ├── init
//   │   └── version
///

// NewServerCommand returns a cobra command for `server` subcommands
func NewServerCommand(cmd *cli.Command) *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "operations regarding the server",
		Long:  "operations affecting the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	serverCmd.AddCommand(
		NewServerConfigCommand(serverCmd),
		NewServerInitCommand(serverCmd),
		NewServerVersionCommand(serverCmd),
	)
	return serverCmd
}

func maybe(err error) {
	if err != nil {
		fmt.Printf("Error:\n%s\n", err)
		os.Exit(1)
	}
}
