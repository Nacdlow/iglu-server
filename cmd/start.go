package cmd

import (
	"log"

	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"
	"net/http"

	"gitlab.com/group-nacdlow/nacdlow-server/routes"
)

// CmdStart represents a command-line command
// which starts the smart home web server.
var CmdStart = cli.Command{
	Name:    "start",
	Aliases: []string{"run"},
	Usage:   "Start the smart home web server",
	Action:  start,
}

func start(clx *cli.Context) (err error) {
	// TODO load configuration files
	// TODO connect to db

	// Start the web server
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	m.Get("/", routes.HomepageHandler)

	log.Fatal(http.ListenAndServe(":8080", m)) // TODO use port from config
	return nil
}
