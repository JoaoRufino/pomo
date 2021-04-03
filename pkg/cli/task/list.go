package task

import (
	"encoding/json"
	"os"
	"time"

	"github.com/joao.rufino/pomo/pkg/cli"
	runnerC "github.com/joao.rufino/pomo/pkg/runner"
	"github.com/joao.rufino/pomo/pkg/server/models"
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
			_list(pomoCli, &options)
		},
	}

	flags := taskListCmd.Flags()

	flags.BoolVarP(&options.asJSON, "json", "j", false, "output task history as JSON")
	flags.BoolVarP(&options.sort, "sort", "s", false, "sort tasks assending in age")
	flags.BoolVarP(&options.all, "all", "a", true, "output all tasks")
	flags.IntVarP(&options.limit, "limit", "n", 0, "limit the number of results by n")
	flags.StringVarP(&options.duration, "duration", "d", "24h", "show tasks within this duration")

	return taskListCmd
}

func _list(pomoCli cli.Cli, options *listOptions) {
	pomoCli.Logger().Debug("Cli request for task list")
	parsed, err := time.ParseDuration(options.duration)
	maybe(err, pomoCli.Logger())

	//get the list from the Server
	list, err := pomoCli.Client().GetTaskList()
	maybe(err, pomoCli.Logger())

	//parse it accordingly
	pomoCli.Logger().Debugf("List has %d tasks %s", len(list))
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
		maybe(json.NewEncoder(os.Stdout).Encode(&list), pomoCli.Logger())
	} else {
		runnerC.SummerizeTasks(pomoCli.Client(), list)
	}
}
