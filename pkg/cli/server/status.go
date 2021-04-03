package server

import (
	"github.com/joao.rufino/pomo/pkg/cli"
	pomo "github.com/joao.rufino/pomo/pkg/server"
	"github.com/spf13/cobra"
)

func NewServerStatusCommand(pomoCli *cli.PomoCli) *cobra.Command {
	serverInitCmd := &cobra.Command{
		Use:   "status",
		Short: "Check server status",
		Long:  `Check server status`,
		Run: func(cmd *cobra.Command, args []string) { // Initialize the databse
			db, err := pomo.NewStore(pomoCli.Config().String("database.path"))
			maybe(err, pomoCli.Logger())
			defer db.Close()
			maybe(pomo.InitDB(db), pomoCli.Logger())
		},
	}
	return serverInitCmd
}
