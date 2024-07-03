package task

import (
	"encoding/json"
	"os"
	"time"

	"github.com/joaorufino/pomo/pkg/cli"
	"github.com/joaorufino/pomo/pkg/core/models"
	runnerC "github.com/joaorufino/pomo/pkg/runner"
	"github.com/spf13/cobra"
)

type listOptions struct {
	asJSON   bool
	sort     bool
	all      bool
	limit    int
	duration string
}

func validateTaskListOptions(opts *listOptions) (*listOptions, error) {

	if opts.limit <= 1 {
		opts.limit = 1
	}

	_, err := time.ParseDuration(opts.duration)
	if err != nil {
		opts.duration = "24h"
	}

	return opts, nil
}

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskListCommand(pomoCli cli.Cli) *cobra.Command {

	options := listOptions{}

	taskListCmd := &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "List tasks",
		Long:  `List all tasks`,
		Run: func(cmd *cobra.Command, args []string) {
			maybe(list(pomoCli, &options), pomoCli.Logger())
		},
	}

	flags := taskListCmd.Flags()

	flags.BoolVarP(&options.asJSON, "json", "j", false, "output task history as JSON")
	flags.BoolVarP(&options.sort, "sort", "s", false, "sort tasks assending in age")
	flags.BoolVarP(&options.all, "all", "a", true, "output all tasks")
	flags.IntVarP(&options.limit, "limit", "n", 0, "limit the number of resultsby n")
	flags.StringVarP(&options.duration, "duration", "d", "24h", "show tasks within this duration")

	return taskListCmd
}

func list(pomoCli cli.Cli, options *listOptions) error {
	pomoCli.Logger().Debug("Cli request for task list")
	parsed, err := time.ParseDuration(options.duration)
	if err != nil {
		return err
	}

	//get the list from the Server

	list := models.List{}
	if plist, err := pomoCli.Client().GetTaskList(); err != nil {
		return err
	} else {
		list = *plist
	}

	//parse it accordingly
	pomoCli.Logger().Debugf("List has %d tasks", len(list))
	if options.sort {
		//sort.Sort(sort.Reverse(list))
	}
	if !options.all {
		list = models.After(time.Now().Add(-parsed), list)
	}
	if options.limit > 0 && (len(list) > options.limit) {
		list = list[0:options.limit]
	}
	if options.asJSON {
		if err := json.NewEncoder(os.Stdout).Encode(&list); err != nil {
			return err
		}
	} else {
		runnerC.SummarizeTasks(pomoCli.Client().Config().String("server.datatimeformat"), list)
	}
	return nil
}
