package server

import (
	"github.com/joao.rufino/pomo/pkg/cli"
	"github.com/joao.rufino/pomo/pkg/server"
	"github.com/joao.rufino/pomo/pkg/store"
	"github.com/spf13/cobra"
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewServerInitCommand(pomoCli *cli.PomoCli) *cobra.Command {
	serverInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Init server",
		Long:  `Start API`,
		Run: func(cmd *cobra.Command, args []string) { // Initialize the databse
			db, err := store.NewStore(pomoCli.Config())
			maybe(err, pomoCli.Logger())
			defer db.Close()
			server, err := server.NewServer(pomoCli.Config(), nil)
			pomoCli.SetServer(&server)
			maybe(err, pomoCli.Logger())
			server.Start()
		},
	}
	return serverInitCmd
}
