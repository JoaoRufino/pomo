package client

import (
	"github.com/joaorufino/pomo/pkg/cli"
	"github.com/joaorufino/pomo/pkg/client"
	"github.com/spf13/cobra"
)

type hostOptions struct {
	port string
	host string
}

// NewClientHostCommand returns a cobra command to host the client locally
func NewClientHostCommand(pomoCli cli.Cli) *cobra.Command {
	options := hostOptions{}

	hostCmd := &cobra.Command{
		Use:   "host",
		Short: "Host a client locally",
		Long:  `Host a client locally`,
		Run: func(cmd *cobra.Command, args []string) {
			startWebServer(&options, pomoCli)
		},
	}

	flags := hostCmd.Flags()
	config := pomoCli.Config()

	// Optional flags
	flags.StringVarP(&options.port, "port", "p", config.Client.HostPort, "Port to host the client")
	flags.StringVarP(&options.host, "host", "H", config.Client.Host, "Host address to bind to")

	return hostCmd
}

func startWebServer(options *hostOptions, pomoCli cli.Cli) {
	pomoCli.Logger().Infof("Starting web server on %s:%s", options.host, options.port)
	client.StartWebServer(pomoCli.Config(), pomoCli.Logger())
	pomoCli.Logger().Infof("Web server started on %s:%s", options.host, options.port)
}
