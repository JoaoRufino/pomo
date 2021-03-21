package cmd

import (
	"database/sql"
	"fmt"
	"time"

	cli "github.com/spf13/cobra"

	"github.com/joao.rufino/pomo/pkg/conf"
	pomo "github.com/joao.rufino/pomo/pkg/internal"
)

// server command
//  pomo
//   ├── server
//   │   ├── config
//   │   ├── init
//   │   └── version
///
func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.AddCommand(configCmd)
	serverCmd.AddCommand(initCmd)
	serverCmd.AddCommand(versionCmd)
}

var (
	serverCmd = &cli.Command{
		Use:   "server",
		Short: "operations regarding the server",
		Long:  "Prints stuff about the user. You could also use the flags in your addPartner() function",
		Run: func(cmd *cli.Command, args []string) {
			fmt.Println("User's name: " + uName)
			fmt.Println("User's number: " + uNumber)
			fmt.Println("User's other stuff: " + uOtherStuff)
		},
	}

	configCmd = &cli.Command{
		Use:   "config",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the databse
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

	versionCmd = &cli.Command{
		Use:   "version",
		Short: "Show version",
		Long:  `Show version`,
		Run: func(cmd *cli.Command, args []string) {
			fmt.Println(conf.Executable + " - " + conf.Version)
		},
	}

	initCmd = &cli.Command{
		Use:   "start",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the databse
			db, err := pomo.newstore(config.dbpath)
			maybe(err)
			defer db.close()
			maybe(pomo.initdb(db))
		},
	}
)
