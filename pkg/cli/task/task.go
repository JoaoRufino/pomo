package task

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
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
func NewTaskCommand(cmd *cli.Command) *cobra.Command {
	taskCmd := &cobra.Command{
		Use:   "task",
		Short: "operations regarding the tasks",
		Long:  "operations affecting the tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	taskCmd.AddCommand(
		NewTaskCreateCommand(taskCmd),
		NewTaskDeleteCommand(taskCmd),
		NewTaskListCommand(taskCmd),
		NewTaskStartCommand(taskCmd),
		NewTaskStatusCommand(taskCmd),
	)
	return taskCmd
}

func maybe(err error) {
	if err != nil {
		fmt.Printf("Error:\n%s\n", err)
		os.Exit(1)
	}
}
