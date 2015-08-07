package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/etcinit/phabulous/app"
	"github.com/facebookgo/inject"
	"github.com/jacobstr/confer"
)

func main() {
	// Create the configuration
	// In this case, we will be using the environment and some safe defaults
	config := confer.NewConfig()
	config.ReadPaths("config/main.yml", "config/main.production.yml")
	config.AutomaticEnv()

	// Next, we setup the dependency graph
	// In this example, the graph won't have many nodes, but on more complex
	// applications it becomes more useful.
	var g inject.Graph
	var phabulous app.Phabulous
	g.Provide(
		&inject.Object{Value: config},
		&inject.Object{Value: &phabulous},
	)
	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Boot the upper layers of the app.
	phabulous.Boot()

	// Setup the command line application
	app := cli.NewApp()
	app.Name = "phabulous"
	app.Usage = "A Phabricator bot in Go"

	// Set version and authorship info
	app.Version = "1.0.0"
	app.Author = "Eduardo Trujillo <ed@chromabits.com>"

	// Setup the default action. This action will be triggered when no
	// subcommand is provided as an argument
	app.Action = func(c *cli.Context) {
		fmt.Println("Usage: phabulous [global options] command [command options] [arguments...]")
	}

	app.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"s", "server", "listen"},
			Usage:   "starts the API server",
			Action:  phabulous.Serve.Run,
		},
	}

	// Begin
	app.Run(os.Args)
}
