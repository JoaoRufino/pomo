package server

import (
	"github.com/joaorufino/pomo/pkg/cli"
	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/server"
	"github.com/joaorufino/pomo/pkg/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewServerInitCommand(pomoCli cli.Cli) *cobra.Command {
	serverInitCmd := &cobra.Command{
		Use:   "run",
		Short: "run server",
		Long:  `Start Server`,
		Run: func(cmd *cobra.Command, args []string) { // Initialize the databse
			db, err := store.NewStore(pomoCli.Config(), pomoCli.Logger())
			maybe(err, pomoCli.Logger())
			defer db.Close()
			serv, err := server.NewServer(pomoCli.Config())
			maybe(err, pomoCli.Logger())
			pomoCli.SetServer(&serv)
			serv.Start()

			conf.Stop.InitInterrupt()
			<-conf.Stop.Chan() // Wait until Stop
			conf.Stop.Wait()   // Wait until everyone cleans up
			_ = zap.L().Sync() // Flush the logger	,
		},
	}
	return serverInitCmd
}
