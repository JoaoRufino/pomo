package task

import (
	"database/sql"

	"github.com/joao.rufino/pomo/pkg/cli"
	runnerC "github.com/joao.rufino/pomo/pkg/runner"
	pomo "github.com/joao.rufino/pomo/pkg/server"
	"github.com/joao.rufino/pomo/pkg/server/comms/unix"
	"github.com/joao.rufino/pomo/pkg/server/models"
	"github.com/spf13/cobra"
)

var (
	taskId *int
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskStartCommand(pomoCli cli.Cli) *cobra.Command {
	taskStartCmd := &cobra.Command{
		Use:   "start",
		Short: "start task",
		Long:  `start a task`,
		Run: func(cmd *cobra.Command, args []string) {
			_start(pomoCli)
		},
	}

	taskId = taskStartCmd.Flags().IntP("taskID", "t", -1, "ID of task to begin")
	taskStartCmd.MarkFlagRequired("taskID")
	return taskStartCmd
}

func _start(pomoCli cli.Cli) {
	db, err := pomo.NewStore(pomoCli.Config().String("database.path"))
	maybe(err, pomoCli.Logger())
	defer db.Close()
	var task *models.Task
	maybe(db.With(func(tx *sql.Tx) error {
		read, err := db.ReadTask(tx, *taskId)
		if err != nil {
			return err
		}
		task = read
		err = db.DeletePomodoros(tx, *taskId)
		if err != nil {
			return err
		}
		task.Pomodoros = []*models.Pomodoro{}
		return nil
	}), pomoCli.Logger())
	runner, err := runnerC.NewTaskRunner(pomoCli, task)
	maybe(err, pomoCli.Logger())
	server, err := unix.NewServer(pomoCli, runner)
	maybe(err, pomoCli.Logger())
	server.Start()
	defer server.Stop()
	runner.Start()
	runnerC.StartUI(runner)
}
