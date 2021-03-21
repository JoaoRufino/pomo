package task

import (
	"database/sql"

	"github.com/joao.rufino/pomo/pkg/conf"
	pomo "github.com/joao.rufino/pomo/pkg/server"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
)

var (
	delete_taskId *int
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskDeleteCommand(cmd *cli.Command) *cobra.Command {
	taskDeleteCmd := &cli.Command{
		Use:   "delete",
		Short: "delete task",
		Long:  `delete task using id`,
		Run: func(cmd *cli.Command, args []string) {
			_delete(args...)
		},
	}

	delete_taskId = taskDeleteCmd.Flags().IntP("taskID", "t", -1, "ID of task to begin")
	taskDeleteCmd.MarkFlagRequired("taskID")
	return taskDeleteCmd
}

func _delete(args ...string) {
	db, err := pomo.NewStore(conf.K.String("database.path"))
	maybe(err)
	defer db.Close()
	maybe(db.With(func(tx *sql.Tx) error {
		return db.DeleteTask(tx, *delete_taskId)
	}))
}
