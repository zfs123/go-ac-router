package acrouter

import (
	"os"

	"github.com/urfave/cli/v2"
)

type CliServer struct {
	App           *cli.App
	apiServer     *ApiServer
	apiCommand    *cli.Command
	apiTlsCommand *cli.Command
	cliCommand    *cli.Command
	docCommand    *cli.Command
	route         *Router
}

func NewCliServer(apiServer *ApiServer, defaultAction func(context *cli.Context) error) *CliServer {
	cliServer := &CliServer{
		App: &cli.App{
			HideHelp: false,
		},
		apiServer: apiServer,
	}
	if defaultAction != nil {
		cliServer.App.Action = defaultAction
	}
	cliServer.setDefaultCommand()
	return cliServer
}

// Set route to generate documentation
func (cs *CliServer) SetRoute(route *Router) {
	cs.route = route
}

// set default command contains cli、doc and api
func (cs *CliServer) setDefaultCommand() {
	cs.apiCommand = &cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "start a api server",
		Action: func(c *cli.Context) error {
			panic(cs.apiServer.Run())
			return nil
		},
	}
	cs.apiTlsCommand = &cli.Command{
		Name:    "tls_server",
		Aliases: []string{"tls"},
		Usage:   "start a api tls server",
		Action: func(c *cli.Context) error {
			panic(cs.apiServer.RunTLS())
			return nil
		},
	}
}

// Add sub-commands of cli mode
func (cs *CliServer) AddCommand(command *cli.Command) {
	cs.App.Commands = append(cs.App.Commands, command)
}

// Run cli app
func (cs *CliServer) Run() error {
	cs.App.Commands = append(cs.App.Commands, cs.apiCommand, cs.apiTlsCommand)
	return cs.App.Run(os.Args)
}
