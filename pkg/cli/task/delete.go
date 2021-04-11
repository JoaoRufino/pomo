package task

import (
	"github.com/joao.rufino/pomo/pkg/cli"
	"github.com/spf13/cobra"
)

type deleteOptions struct {
	taskID int
}

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskDeleteCommand(pomoCli cli.Cli) *cobra.Command {

	options := &deleteOptions{}

	taskDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete task",
		Long:  `delete task using id`,
		Run: func(cmd *cobra.Command, args []string) {
			delete(pomoCli, options)
		},
	}

	flags := taskDeleteCmd.Flags()

	flags.IntVarP(&options.taskID, "taskID", "t", -1, "ID of task to begin")
	taskDeleteCmd.MarkFlagRequired("taskID")

	return taskDeleteCmd
}

func delete(pomoCli cli.Cli, options *deleteOptions) {
	err := pomoCli.Client().DeleteTaskByID(options.taskID)
	maybe(err, pomoCli.Logger())
}
