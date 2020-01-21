package cmd

import (
	"fmt"
	"log"

	"net/http"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
	"gitlab.com/group-nacdlow/nacdlow-server/routes"
)

// CmdStart represents a command-line command
// which starts the smart home web server.
var CmdStart = cli.Command{
	Name:    "start",
	Aliases: []string{"run"},
	Usage:   "Start the smart home web server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "port",
			Value: "8080",
			Usage: "the web server port",
		},
	},
	Action: start,
}

func start(clx *cli.Context) (err error) {
	settings.LoadConfig()
	engine := models.SetupEngine()
	defer engine.Close()

	// Start the web server
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())

	m.NotFound(routes.NotFoundHandler)

	m.Get("/", routes.HomepageHandler)
	m.Get("/dashboard", routes.DashboardHandler)
	m.Get("/devices", routes.DevicesHandler)
	m.Get("/lights", routes.LightsHandler)
	m.Get("/temperature", routes.HeatingHandler)
	m.Get("/speakers", routes.SpeakerHandler)
	m.Group("/room", func() {
		m.Get("/add", routes.AddRoomHandler)
		m.Get("/:name", routes.SpecificRoomsHandler)
	})
	m.Get("/rooms", routes.RoomsHandler)
	m.Get("/register", routes.RegisterHandler)

	log.Printf("Starting server on port %s!\n", clx.String("port"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", clx.String("port")), m))
	return nil
}
