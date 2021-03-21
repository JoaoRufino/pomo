package task

import (
	"database/sql"
	"time"

	"github.com/joao.rufino/pomo/pkg/conf"
	runnerC "github.com/joao.rufino/pomo/pkg/runner"
	pomo "github.com/joao.rufino/pomo/pkg/server"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
)

var (
	duration  *string
	message   *string
	pomodoros *int
	start     *bool
	tags      *[]string
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskCreateCommand(cmd *cli.Command) *cobra.Command {
	taskCreateCmd := &cli.Command{
		Use:   "create",
		Short: "create task",
		Long:  `create task`,
		Run: func(cmd *cli.Command, args []string) {
			_create(args...)
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
func _create(args ...string) {
	parsed, err := time.ParseDuration(*duration)
	maybe(err)
	db, err := pomo.NewStore(conf.K.String("database.path"))
	maybe(err)
	defer db.Close()
	task := &pomo.Task{
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
	}))
	if *start {
		runner, err := runnerC.NewTaskRunner(task)
		maybe(err)
		server, err := pomo.NewServer(runner)
		maybe(err)
		server.Start()
		defer server.Stop()
		runner.Start()
		runnerC.StartUI(runner)
	}
}
