package server

import (
	"github.com/joao.rufino/pomo/pkg/cli"
	"github.com/joao.rufino/pomo/pkg/store"
	"github.com/spf13/cobra"
)

func NewServerStatusCommand(pomoCli *cli.PomoCli) *cobra.Command {
	serverInitCmd := &cobra.Command{
		Use:   "status",
		Short: "Check server status",
		Long:  `Check server status`,
		Run: func(cmd *cobra.Command, args []string) { // Initialize the databse
			db, err := store.NewStore(pomoCli.Config())
			maybe(err, pomoCli.Logger())
			defer db.Close()
			maybe(db.InitDB(), pomoCli.Logger())
		},
	}
	return serverInitCmd
}
