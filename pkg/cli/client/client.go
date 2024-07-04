package client

import (
	"os"

	"github.com/joaorufino/pomo/pkg/cli"
	"github.com/joaorufino/pomo/pkg/client"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewClientCommand(pomoCli cli.Cli) *cobra.Command {
	clientCmd := &cobra.Command{
		Use:   "client",
		Short: "operations regarding the clients",
		Long:  "operations affecting the clients",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient(pomoCli.Config())
			maybe(err, pomoCli.Logger())
			pomoCli.SetClient(&c)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			pomoCli.Client().Close()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	clientCmd.AddCommand(
		NewClientHostCommand(pomoCli),
	)
	return clientCmd
}

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:%s\n", err)
		os.Exit(1)
	}
}
