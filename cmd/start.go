package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"net/http"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/urfave/cli/v2"
	macaron "gopkg.in/macaron.v1"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/models/forms"
	forms_sim "gitlab.com/group-nacdlow/nacdlow-server/models/forms/sim"
	"gitlab.com/group-nacdlow/nacdlow-server/models/simulation"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
	"gitlab.com/group-nacdlow/nacdlow-server/routes"
	routes_sim "gitlab.com/group-nacdlow/nacdlow-server/routes/simulator"
)

// CmdStart represents a command-line command
// which starts the smart home web server.
var CmdStart = &cli.Command{
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
	go simulation.Start()

	// Start the web server
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())

	m.Use(routes.ContextInit())

	m.NotFound(routes.NotFoundHandler)

	m.Get("/", routes.LoginHandler)
	m.Post("/", binding.Bind(forms.SignInForm{}), routes.PostLoginHandler)
	m.Get("/register", routes.RegisterHandler)
	m.Post("/register", routes.PostRegisterHandler)

	m.Group("", func() {
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
	}, routes.RequireLogin)

	m.Group("/sim", func() {
		m.Get("/", routes_sim.HomepageHandler)
		m.Get("/data.json", routes_sim.DataHandler)
		m.Post("/override_weather", binding.Bind(forms_sim.OverrideWeatherForm{}),
			routes_sim.PostOverrideWeatherHandler)
		m.Post("/time_sleep", binding.Bind(forms_sim.ChangeTimeSleepForm{}),
			routes_sim.PostChangeTimeSleepHandler)
		m.Get("/env_status", routes_sim.EnvStatusHandler)
		m.Get("/toggle/:id", routes_sim.ToggleHandler)
	})

	log.WithFields(log.Fields{"port": clx.String("port")}).Printf("Starting server.")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", clx.String("port")), m))
	return nil
}
