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
			_status(pomoCli)
		},
	}

	return taskStatusCmd
}

func _status(pomoCli cli.Cli) {
	status, err := pomoCli.Client().GetServerStatus()
	maybe(err, pomoCli.Logger())
	runnerC.OutputStatus(*status)
}
