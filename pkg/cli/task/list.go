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

var (
	asJSON *bool
	assend *bool
	all    *bool
	limit  *int
	period *string
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskListCommand(cmd *cli.Command) *cobra.Command {
	taskListCmd := &cli.Command{
		Use:   "list",
		Short: "List tasks",
		Long:  `List all tasks`,
		Run: func(cmd *cli.Command, args []string) {
			_list(args...)
		},
	}

	asJSON = taskListCmd.Flags().BoolP("json", "j", false, "output task history as JSON")
	assend = taskListCmd.Flags().BoolP("sort", "s", false, "sort tasks assending in age")
	all = taskListCmd.Flags().BoolP("all", "a", true, "output all tasks")
	limit = taskListCmd.Flags().IntP("limit", "n", 0, "limit the number of results by n")
	period = taskListCmd.Flags().StringP("duration", "d", "24h", "show tasks within this duration")

	return taskListCmd
}

func _list(args ...string) {
	parsed, err := time.ParseDuration(*period)
	maybe(err)
	db, err := pomo.NewStore(conf.K.String("database.path"))
	maybe(err)
	defer db.Close()
	maybe(db.With(func(tx *sql.Tx) error {
		tasks, err := db.ReadTasks(tx)
		maybe(err)
		if *assend {
			sort.Sort(sort.Reverse(pomo.ByID(tasks)))
		}
		if !*all {
			tasks = pomo.After(time.Now().Add(-parsed), tasks)
		}
		if *limit > 0 && (len(tasks) > *limit) {
			tasks = tasks[0:*limit]
		}
		if *asJSON {
			maybe(json.NewEncoder(os.Stdout).Encode(tasks))
			return nil
		}
		maybe(err)
		runnerC.SummerizeTasks(tasks)
		return nil
	}))
}
