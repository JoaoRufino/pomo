package task

import (
	"github.com/joao.rufino/pomo/pkg/cli"
	"github.com/spf13/cobra"

	runnerC "github.com/joao.rufino/pomo/pkg/runner"
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskStatusCommand(pomoCli cli.Cli) *cobra.Command {
	taskStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "task status",
		Long:  `request the tasks status`,
		Run: func(cmd *cobra.Command, args []string) {
			maybe(status(pomoCli), pomoCli.Logger())
		},
	}

	return taskStatusCmd
}

func status(pomoCli cli.Cli) error {
	status, err := pomoCli.Client().GetServerStatus()
	if err != nil {
		return err
	}
	runnerC.OutputStatus(*status)
	return nil
}
