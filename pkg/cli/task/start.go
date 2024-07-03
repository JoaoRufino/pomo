package task

import (
	"errors"

	"github.com/joao.rufino/pomo/pkg/cli"
	"github.com/spf13/cobra"
)

type startOptions struct {
	taskID int
}

// NewStartCommand returns a cobra command for `config` subcommands
func NewTaskStartCommand(pomoCli cli.Cli) *cobra.Command {

	options := startOptions{}

	taskStartCmd := &cobra.Command{
		Use:   "start",
		Short: "start task",
		Long:  `start a task`,
		Run: func(cmd *cobra.Command, args []string) {
			maybe(start(pomoCli, &options), pomoCli.Logger())
		},
	}

	flags := taskStartCmd.Flags()

	flags.IntVarP(&options.taskID, "taskID", "t", -1, "ID of task to begin")
	taskStartCmd.MarkFlagRequired("taskID")

	return taskStartCmd
}

func start(pomoCli cli.Cli, options *startOptions) error {
	if pomoCli.Client() == nil {
		return errors.New("client not defined")
	}
	pomoCli.Client().StartTask(options.taskID)
	return nil
}
