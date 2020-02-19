package cmd

import (
	"github.com/go-macaron/binding"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/urfave/cli/v2"
	macaron "gopkg.in/macaron.v1"

	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"time"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/models/forms"
	forms_sim "gitlab.com/group-nacdlow/nacdlow-server/models/forms/sim"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/plugin"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/simulation"
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

func getMacaron() *macaron.Macaron {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())
	m.Use(macaron.Static("public",
		macaron.StaticOptions{
			Expires: func() string {
				return time.Now().Add(24 * 60 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
			},
		}))

	m.Use(routes.ContextInit())

	// Load plugin middlewares
	for _, pl := range plugin.LoadedPlugins {
		if pl.Middleware != nil {
			m.Use(pl.Middleware())
		}
	}

	m.NotFound(routes.NotFoundHandler)

	m.Get("/", routes.LoginHandler)
	m.Get("/login", routes.LoginHandler)
	m.Post("/", binding.Bind(forms.SignInForm{}), routes.PostLoginHandler)
	m.Get("/register", routes.RegisterHandler)
	//m.Post("/register", routes.PostRegisterHandler)
	m.Post("/register", binding.Bind(forms.RegisterForm{}), routes.AddUserHandler) //registers a user

	m.Group("", func() {
		m.Get("/dashboard", routes.DashboardHandler)
		m.Get("/logout", routes.LogoutHandler)
		m.Group("/room", func() {
			m.Group("/:name", func() {
				m.Get("", routes.SpecificRoomsHandler)
			})
		})
		m.Get("/overview", routes.OverviewHandler)
		m.Get("/rooms", routes.RoomsHandler)
		m.Get("/toggle_device/:id", routes.ToggleHandler)
		m.Get("/toggle_slider/:id/:value", routes.SliderHandler)
		m.Get("/alerts", routes.AlertsHandler)

		m.Get("/toggle_fave/:id", routes.FaveHandler)     //set device as fave
		m.Get("/remove_device/:id", routes.RemoveHandler) //remove a device

		m.Group("", func() {
			m.Group("/add", func() {
				m.Get("", routes.AddHandler)
				m.Get("/room", routes.AddRoomHandler)
				m.Post("/room", binding.Bind(forms.AddRoomForm{}),
					routes.PostRoomHandler)
				m.Get("/device", routes.AddDeviceHandler)
				m.Post("/device", binding.Bind(forms.AddDeviceForm{}),
					routes.AddDeviceRoomPostHandler)
			})
			m.Get("/settings", routes.SettingsHandler)
			m.Group("/settings/plugins", func() {
				m.Get("", routes.PluginsSettingsHandler)
				m.Get("/:id", routes.InstallPluginSettingsHandler)
				m.Get("/confirm/:id", routes.InstallPluginConfirmSettingsHandler) // TODO use POST so it is secure
			})
			m.Get("/settings/accounts", routes.AccountSettingsHandler)
			m.Post("/settings/accounts", routes.PostAccountSettingsHandler)
			m.Get("/settings/appearance", routes.AppearanceSettingsHandler)
		}, routes.RequireAdmin)

		m.Get("/battery_stat", routes.BatteryStatHandler)
	}, routes.RequireLogin)

	// Simulator routes (no auth)
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

	// For debugging purposes.
	m.Group("/debug/pprof", func() {
		m.Get("/", pprofHandler(pprof.Index))
		m.Get("/cmdline", pprofHandler(pprof.Cmdline))
		m.Get("/profile", pprofHandler(pprof.Profile))
		m.Post("/symbol", pprofHandler(pprof.Symbol))
		m.Get("/symbol", pprofHandler(pprof.Symbol))
		m.Get("/trace", pprofHandler(pprof.Trace))
		m.Get("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
		m.Get("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		m.Get("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		m.Get("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		m.Get("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
		m.Get("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
	})
	return m
}

func pprofHandler(h http.HandlerFunc) macaron.Handler {
	handler := http.HandlerFunc(h)
	return func(c *macaron.Context) {
		handler.ServeHTTP(c.Resp, c.Req.Request)
	}
}

func start(clx *cli.Context) (err error) {
	// Log both to a file and to stdout
	file, err := os.OpenFile("iglu.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	multi := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multi)

	settings.LoadConfig()
	engine := models.SetupEngine()
	defer engine.Close()
	go simulation.Start()
	plugin.LoadPlugins()

	// Start the web server
	m := getMacaron()

	log.Printf("Starting TLS web server on :%s\n", clx.String("port"))
	server := &http.Server{Addr: fmt.Sprintf(":%s", clx.String("port")), Handler: m}
	go func() {
		log.Fatal(server.ListenAndServeTLS("fullchain.pem", "privkey.pem"))
	}()

	// Capture system interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	return nil
}
