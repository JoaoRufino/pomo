package task

import (
	"os"

	"github.com/joao.rufino/pomo/pkg/cli"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// task command
//  pomo
//   ├── task
//   │   ├── create
//   │   ├── delete
//   │   ├── list
//   │   ├── start
//   │   └── status
///
// NewServerCommand returns a cobra command for `server` subcommands
func NewTaskCommand(pomoCli cli.Cli) *cobra.Command {
	taskCmd := &cobra.Command{
		Use:   "task",
		Short: "operations regarding the tasks",
		Long:  "operations affecting the tasks",
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
