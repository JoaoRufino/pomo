package task

import (
	"github.com/joao.rufino/pomo/pkg/cli"
	"github.com/joao.rufino/pomo/pkg/runner/client"
	"github.com/joao.rufino/pomo/pkg/server/models"
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
	client, err := client.NewClient(pomoCli.Config())
	if err != nil {
		runnerC.OutputStatus(models.Status{})
		return
	}
	defer client.Close()
	status, err := client.Status()
	maybe(err, pomoCli.Logger())
	runnerC.OutputStatus(*status)
}
