
package cmd

import (
	"database/sql"
	"time"

	cli "github.com/spf13/cobra"
"github.com/joao.rufino/pomo/pkg/conf"
	pomo "github.com/joao.rufino/pomo/pkg/internal"
)

// task command
//  pomo
//   ├── task
//   │   ├── begin
//   │   ├── create
//   │   ├── delete
//   │   ├── list
//   │   ├── start
//   │   └── status
///
func init() {
	rootCmd.AddCommand(taskCmd)
	
	taskCmd.AddCommand(beginCmd)
	taskCmd.AddCommand(createCmd)
	taskCmd.AddCommand(deleteCmd)
	taskCmd.AddCommand(listCmd)
	taskCmd.AddCommand(startCmd)
	taskCmd.AddCommand(statusCmd)
}

var (
    taskCmd = &cli.Command{
		Use:   "status st",
		Short: "Output the current status",
		Long:  `Output the current status`,
		Run: func(cmd *cli.Command, args []string) {
		cmd.Action = func() {}
	},
}
	beginCmd = &cli.Command{
		Use:   "start",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) {
			taskID, _ := cmd.Flags().GetString("taskID")
			message =
			duration =
			tags =
			pomodoros =
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := pomo.NewStore(config.K.String("database.path"))
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
			runner, err := pomo.NewTaskRunner(task)
			maybe(err)
			server, err := pomo.NewServer(runner)
			maybe(err)
			server.Start()
			defer server.Stop()
			runner.Start()
			pomo.StartUI(runner)
		},
	}
	createCmd = &cli.Command{
		Use:   "start",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the databse
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := pomo.NewStore(config.DBPath)
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
			runner, err := pomo.NewTaskRunner(task, config)
			maybe(err)
			server, err := pomo.NewServer(runner, config)
			maybe(err)
			server.Start()
			defer server.Stop()
			runner.Start()
			pomo.StartUI(runner)
		},
	}

	deleteCmd = &cli.Command{
		Use:   "start",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the databse
			cmd.Spec = "[OPTIONS] TASK_ID"
			var taskID = cmd.IntArg("TASK_ID", -1, "task to delete")
			cmd.Action = func() {
				db, err := pomo.NewStore(config.DBPath)
				maybe(err)
				defer db.Close()
				maybe(db.With(func(tx *sql.Tx) error {
					return db.DeleteTask(tx, *taskID)
				}))
			}
		},
	}

	listCmd = &cli.Command{
		Use:   "start",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the databse
			duration, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := pomo.NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(tx *sql.Tx) error {
				tasks, err := db.ReadTasks(tx)
				maybe(err)
				if *assend {
					sort.Sort(sort.Reverse(pomo.ByID(tasks)))
				}
				if !*all {
					tasks = pomo.After(time.Now().Add(-duration), tasks)
				}
				if *limit > 0 && (len(tasks) > *limit) {
					tasks = tasks[0:*limit]
				}
				if *asJSON {
					maybe(json.NewEncoder(os.Stdout).Encode(tasks))
					return nil
				}
				maybe(err)
				pomo.SummerizeTasks(config, tasks)
				return nil
			}))
		},
	}

	startCmd = &cli.Command{
		Use:   "start",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) {
			taskID, _ := cmd.Flags().GetString("taskID")
			message =
			duration =
			tags =
			pomodoros =
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := pomo.NewStore(config.K.String("database.path"))
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
			runner, err := pomo.NewTaskRunner(task)
			maybe(err)
			server, err := pomo.NewServer(runner)
			maybe(err)
			server.Start()
			defer server.Stop()
			runner.Start()
			pomo.StartUI(runner)
		},
	}

	statusCmd = &cli.Command{
		Use:   "status st",
		Short: "Output the current status",
		Long:  `Output the current status`,
		Run: func(cmd *cli.Command, args []string) {
		cmd.Action = func() {
			client, err := pomo.NewClient(config.C.String("server.socket"))
			if err != nil {
				pomo.OutputStatus(pomo.Status{})
				return
			}
			defer client.Close()
			status, err := client.Status()
			maybe(err)
			pomo.OutputStatus(*status)	
		}
	},
}
	
)
