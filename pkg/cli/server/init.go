package server

import (
	"github.com/joao.rufino/pomo/pkg/conf"
	pomo "github.com/joao.rufino/pomo/pkg/server"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewServerInitCommand(cmd *cli.Command) *cobra.Command {
	serverInitCmd := &cli.Command{
		Use:   "init",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the databse
			db, err := pomo.NewStore(conf.K.String("database.path"))
			maybe(err)
			defer db.Close()
			maybe(pomo.InitDB(db))
		},
	}
	return serverInitCmd
}
