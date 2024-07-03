package root

import (
	"fmt"
	"log"
	"os"

	_ "net/http/pprof" // Import for pprof

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/joaorufino/pomo/pkg/cli"
	"github.com/joaorufino/pomo/pkg/cli/server"
	"github.com/joaorufino/pomo/pkg/cli/task"
	"github.com/joaorufino/pomo/pkg/conf"
)

var (
	// Config and global logger
	pidFile string
)

func NewRootCommand(pomoCli *cli.PomoCli) *cobra.Command {
	return &cobra.Command{
		Version:           pomoCli.Version(),
		Use:               pomoCli.Executable(),
		PersistentPreRunE: prerun,
		PersistentPostRun: cleanup,
	}

}

// Execute starts the program
func Execute() {
	// load default config
	pomoCli, err := cli.NewPomoCli("")
	if err != nil {
		log.Fatal(err)
	}
	rootCmd := NewRootCommand(pomoCli)
	configFile := rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
	if configFile != nil && *configFile != "" {
		config, err := conf.LoadConfig(*configFile)
		maybe(err, pomoCli.Logger())
		pomoCli.SetConfig(config)
	}
	rootCmd.AddCommand(
		server.NewServerCommand(pomoCli),
		task.NewTaskCommand(pomoCli))

	// Run the program
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func prerun(cmd *cobra.Command, args []string) error {
	//keep track of processID this will make sure we can keep track on unexpected behaviour
	//logic adapted from "Go Systems Programming - Milhalis Tsoukalos"
	//"https://man7.org/linux/man-pages/man2/open.2.html"
	// Create Pid File
	pidFile = "" //TODO placeholder
	if pidFile != "" {
		file, err := os.OpenFile(pidFile,
			os.O_CREATE| //create if it doesnt exist
				os.O_TRUNC| //truncates file if it already exists
				os.O_WRONLY, //write only
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
