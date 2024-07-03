package task

import (
	"time"

	"github.com/joaorufino/pomo/pkg/cli"
	"github.com/joaorufino/pomo/pkg/core/models"
	"github.com/spf13/cobra"
)

type createOptions struct {
	duration  string
	message   string
	pomodoros int
	start     bool
	tags      []string
}

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskCreateCommand(pomoCli cli.Cli) *cobra.Command {

	options := createOptions{}

	taskCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "create task",
		Long:  `create task`,
		Run: func(cmd *cobra.Command, args []string) {
			create(pomoCli, &options)
		},
	}

	flags := taskCreateCmd.Flags()

	//optional flags
	flags.StringVarP(&options.duration, "duration", "d", "25m", "duration of each stent")
	flags.IntVarP(&options.pomodoros, "pomodoros", "p", 4, "number of pomodoros")
	flags.StringSliceVarP(&options.tags, "tag", "t", []string{}, "tags associated with this task")
	flags.BoolVarP(&options.start, "start", "s", false, "start pomodoro after creation")

	//mandatory flags
	flags.StringVarP(&options.message, "message", "m", "", "descriptive name of the given task")
	taskCreateCmd.MarkFlagRequired("message")

	return taskCreateCmd
}

// creates a task
func create(pomoCli cli.Cli, options *createOptions) {
	parsed, err := time.ParseDuration(options.duration)
	maybe(err, pomoCli.Logger())

	task := &models.Task{
		Message:    options.message,
		Tags:       options.tags,
		NPomodoros: options.pomodoros,
		Duration:   parsed,
	}
	taskID, err := pomoCli.Client().CreateTask(task)
	maybe(err, pomoCli.Logger())

	//if the user requested to start the created task
	if options.start {
		pomoCli.Client().StartTask(taskID)
	} else {
		pomoCli.Logger().Debugf("Task id: %d created", taskID)
	}
}
