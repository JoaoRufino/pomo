package root

import (
	"fmt"
	"log"
	"os"

	_ "net/http/pprof" // Import for pprof

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/joaorufino/pomo/pkg/cli"
	"github.com/joaorufino/pomo/pkg/cli/client"
	"github.com/joaorufino/pomo/pkg/cli/server"
	"github.com/joaorufino/pomo/pkg/cli/task"
	"github.com/joaorufino/pomo/pkg/conf"
)

var (
	// Config and global logger
	pidFile    string
	configFile string
)

// NewRootCommand initializes the root command
func NewRootCommand(pomoCli *cli.PomoCli) *cobra.Command {
	rootCmd := &cobra.Command{
		Version:           pomoCli.Version(),
		Use:               pomoCli.Executable(),
		PersistentPreRunE: prerun,
		PersistentPostRun: cleanup,
	}

	// Define persistent flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	return rootCmd
}

// Execute starts the program
func Execute() {
	pomoCli, err := cli.NewPomoCli("")
	if err != nil {
		log.Fatal(err)
	}
	rootCmd := NewRootCommand(pomoCli)

	// Parse the flags early to load the config file
	rootCmd.PersistentFlags().Parse(os.Args[1:])
	if configFile != "" {
		config, err := conf.LoadConfig(configFile)
		maybe(err, pomoCli.Logger())
		pomoCli.SetConfig(config)
	}

	rootCmd.AddCommand(
		server.NewServerCommand(pomoCli),
		task.NewTaskCommand(pomoCli),
		client.NewClientCommand(pomoCli),
	)

	// Run the program
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func prerun(cmd *cobra.Command, args []string) error {
	// Create PID File
	if pidFile != "" {
		file, err := os.OpenFile(pidFile,
			os.O_CREATE| // create if it doesn't exist
				os.O_TRUNC| // truncates file if it already exists
				os.O_WRONLY, // write only
			0666)
		if err != nil {
			return fmt.Errorf("could not create pid file: %s Error:%v", pidFile, err)
		}
		defer file.Close()
		_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
		if err != nil {
			return fmt.Errorf("could not create pid file: %s Error:%v", pidFile, err)
		}
	}
	return nil
}

func cleanup(cmd *cobra.Command, args []string) {
	// PID Cleanup
	if pidFile != "" {
		os.Remove(pidFile)
	}
}

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:\n%s\n", err)
		os.Exit(1)
	}
}
