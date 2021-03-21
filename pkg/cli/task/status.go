package task

import (
	"github.com/joao.rufino/pomo/pkg/conf"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"

	runnerC "github.com/joao.rufino/pomo/pkg/runner"
	pomo "github.com/joao.rufino/pomo/pkg/server"
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskStatusCommand(cmd *cli.Command) *cobra.Command {
	taskStatusCmd := &cli.Command{
		Use:   "status",
		Short: "task status",
		Long:  `request the tasks status`,
		Run: func(cmd *cli.Command, args []string) {
			_status(args...)
		},
	}

	return taskStatusCmd
}

func _status(args ...string) {
	client, err := pomo.NewClient(conf.K.String("server.socket"))
	if err != nil {
		runnerC.OutputStatus(pomo.Status{})
		return
	}
	defer client.Close()
	status, err := client.Status()
	maybe(err)
	runnerC.OutputStatus(*status)
}
