package cmd

import (
	"fmt"
	"os"

	_ "net/http/pprof" // Import for pprof

	cli "github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/joao.rufino/pomo/pkg/conf"
	pomo "github.com/joao.rufino/pomo/pkg/internal"
)

var (

	// Config and global logger
	pidFile string
	logger  *zap.SugaredLogger
	config  *pomo.Config

	// The Root Cli Handler
	rootCmd = &cli.Command{
		Version:           conf.GitVersion,
		Use:               conf.Executable,
		PersistentPreRunE: prerun,
		PersistentPostRun: cleanup,
	}
)

// Execute starts the program
func Execute() {

	// Load configuration
	_ = conf.Defaults(conf.K)
	configFile := rootCmd.PersistentFlags().StringP("config", "c", "", "setup the config file")

	if configFile != nil && *configFile != "" {
		_ = conf.File(conf.K, *configFile)
	}
	_ = conf.Env(conf.K)

	conf.InitLogger(conf.K)

	logger = zap.S().With("package", "cmd")

	// Run the program
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func prerun(cmd *cli.Command, args []string) error {
	//keep track of processID this will make sure we can keep track on unexpected behaviour
	//logic adapted from "Go Systems Programming - Milhalis Tsoukalos"
	//"https://man7.org/linux/man-pages/man2/open.2.html"
	// Create Pid File
	pidFile = conf.C.String("pidfile")
	if pidFile != "" {
		file, err := os.OpenFile(pidFile,
			os.O_CREATE| //create if it doesnt exist
				os.O_TRUNC| //truncates file if it already exists
				os.O_WRONLY, //write only
			0666)
		if err != nil {
			return fmt.Errorf("Could not create pid file: %s Error:%v", pidFile, err)
		}
		defer file.Close()
		_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
		if err != nil {
			return fmt.Errorf("Could not create pid file: %s Error:%v", pidFile, err)
		}
	}
	return nil
}

func cleanup(cmd *cli.Command, args []string) {
	// PID Cleanup
	if pidFile != "" {
		os.Remove(pidFile)
	}
}

func maybe(err error) {
	if err != nil {
		fmt.Printf("Error:\n%s\n", err)
		os.Exit(1)
	}
}
