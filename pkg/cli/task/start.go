package task

import (
	"database/sql"

	"github.com/joao.rufino/pomo/pkg/conf"
	runnerC "github.com/joao.rufino/pomo/pkg/runner"
	pomo "github.com/joao.rufino/pomo/pkg/server"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
)

var (
	taskId *int
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskStartCommand(cmd *cli.Command) *cobra.Command {
	taskStartCmd := &cli.Command{
		Use:   "start",
		Short: "start task",
		Long:  `start a task`,
		Run: func(cmd *cli.Command, args []string) {
			_start(args...)
		},
	}

	taskId = taskStartCmd.Flags().IntP("taskID", "t", -1, "ID of task to begin")
	taskStartCmd.MarkFlagRequired("taskID")
	return taskStartCmd
}

func _start(args ...string) {
	db, err := pomo.NewStore(conf.K.String("database.path"))
	maybe(err)
	defer db.Close()
	var task *pomo.Task
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
		task.Pomodoros = []*pomo.Pomodoro{}
		return nil
	}))
	runner, err := runnerC.NewTaskRunner(task)
	maybe(err)
	server, err := pomo.NewServer(runner)
	maybe(err)
	server.Start()
	defer server.Stop()
	runner.Start()
	runnerC.StartUI(runner)
}
