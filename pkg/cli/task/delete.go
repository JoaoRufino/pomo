package task

import (
	"database/sql"

	"github.com/joao.rufino/pomo/pkg/cli"
	pomo "github.com/joao.rufino/pomo/pkg/server"
	"github.com/spf13/cobra"
)

var (
	delete_taskId *int
)

// NewConfigCommand returns a cobra command for `config` subcommands
func NewTaskDeleteCommand(pomoCli cli.Cli) *cobra.Command {
	taskDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete task",
		Long:  `delete task using id`,
		Run: func(cmd *cobra.Command, args []string) {
			_delete(pomoCli)
		},
	}

	delete_taskId = taskDeleteCmd.Flags().IntP("taskID", "t", -1, "ID of task to begin")
	taskDeleteCmd.MarkFlagRequired("taskID")
	return taskDeleteCmd
}

func _delete(pomoCli cli.Cli) {
	db, err := pomo.NewStore(pomoCli.Config().String("database.path"))
	maybe(err, pomoCli.Logger())
	defer db.Close()
	maybe(db.With(func(tx *sql.Tx) error {
		return db.DeleteTask(tx, *delete_taskId)
	}), pomoCli.Logger())
}
