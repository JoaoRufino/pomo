package task

import (
	"database/sql"
	"encoding/json"
	"os"
	"sort"
	"time"

	"github.com/joao.rufino/pomo/pkg/conf"
	runnerC "github.com/joao.rufino/pomo/pkg/runner"
	pomo "github.com/joao.rufino/pomo/pkg/server"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
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
func NewTaskListCommand(cmd *cli.Command) *cobra.Command {

	options := listOptions{}

	taskListCmd := &cli.Command{
		Use:   "list [OPTIONS]",
		Short: "List tasks",
		Long:  `List all tasks`,
		Run: func(cmd *cli.Command, args []string) {
			_list(&options)
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

func _list(options *listOptions) {
	parsed, err := time.ParseDuration(options.duration)
	maybe(err)
	db, err := pomo.NewStore(conf.K.String("database.path"))
	maybe(err)
	defer db.Close()
	maybe(db.With(func(tx *sql.Tx) error {
		tasks, err := db.ReadTasks(tx)
		maybe(err)
		if options.sort {
			sort.Sort(sort.Reverse(pomo.ByID(tasks)))
		}
		if !options.all {
			tasks = pomo.After(time.Now().Add(-parsed), tasks)
		}
		if options.limit > 0 && (len(tasks) > options.limit) {
			tasks = tasks[0:options.limit]
		}
		if options.asJSON {
			maybe(json.NewEncoder(os.Stdout).Encode(tasks))
			return nil
		}
		maybe(err)
		runnerC.SummerizeTasks(tasks)
		return nil
	}))
}
