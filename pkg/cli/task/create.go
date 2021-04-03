package task

import (
	"database/sql"
	"time"

	"github.com/joao.rufino/pomo/pkg/cli"
	runnerC "github.com/joao.rufino/pomo/pkg/runner"
	pomo "github.com/joao.rufino/pomo/pkg/server"
	"github.com/joao.rufino/pomo/pkg/server/comms/unix"
	"github.com/joao.rufino/pomo/pkg/server/models"
	"github.com/spf13/cobra"
)

var (
	duration  *string
	message   *string
	pomodoros *int
	start     *bool
	tags      *[]string
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskCreateCommand(pomoCli cli.Cli) *cobra.Command {
	taskCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "create task",
		Long:  `create task`,
		Run: func(cmd *cobra.Command, args []string) {
			_create(pomoCli)
		},
	}
	//optional flags
	duration = taskCreateCmd.Flags().StringP("duration", "d", "25m", "duration of each stent")
	pomodoros = taskCreateCmd.Flags().IntP("pomodoros", "p", 4, "number of pomodoros")
	tags = taskCreateCmd.Flags().StringSliceP("tag", "t", []string{}, "tags associated with this task")
	start = taskCreateCmd.Flags().BoolP("start", "s", false, "start pomodoro after creation")

	//mandatory flags
	message = taskCreateCmd.Flags().StringP("message", "m", "", "descriptive name of the given task")
	taskCreateCmd.MarkFlagRequired("message")

	return taskCreateCmd
}

// creates a task
func _create(pomoCli cli.Cli) {
	parsed, err := time.ParseDuration(*duration)
	maybe(err, pomoCli.Logger())
	db, err := pomo.NewStore(pomoCli.Config().String("database.path"))
	maybe(err, pomoCli.Logger())
	defer db.Close()
	task := &models.Task{
		Message:    *message,
		Tags:       *tags,
		NPomodoros: *pomodoros,
		Duration:   parsed,
	}
	maybe(db.With(func(tx *sql.Tx) error {
		id, err := db.CreateTask(tx, *task)
		if err != nil {
			return err
		}
		task.ID = id
		return nil
	}), pomoCli.Logger())
	if *start {
		runner, err := runnerC.NewTaskRunner(pomoCli, task)
		maybe(err, pomoCli.Logger())
		server, err := unix.NewServer(pomoCli, runner)
		maybe(err, pomoCli.Logger())
		server.Start()
		defer server.Stop()
		runner.Start()
		runnerC.StartUI(runner)
	}
}
