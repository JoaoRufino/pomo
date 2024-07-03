package task

import (
	"os"

	"github.com/joaorufino/pomo/pkg/cli"
	"github.com/joaorufino/pomo/pkg/client"
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	c core.Client
)

// task command
//
//	pomo
//	 ├── task
//	 │   ├── create
//	 │   ├── delete
//	 │   ├── list
//	 │   ├── start
//	 │   └── status
//
// /
// NewServerCommand returns a cobra command for `server` subcommands
func NewTaskCommand(pomoCli cli.Cli) *cobra.Command {
	taskCmd := &cobra.Command{
		Use:   "task",
		Short: "operations regarding the tasks",
		Long:  "operations affecting the tasks",
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
	taskCmd.AddCommand(
		NewTaskCreateCommand(pomoCli),
		NewTaskDeleteCommand(pomoCli),
		NewTaskListCommand(pomoCli),
		NewTaskStartCommand(pomoCli),
		NewTaskStatusCommand(pomoCli),
	)
	return taskCmd
}

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:%s\n", err)
		os.Exit(1)
	}
}
