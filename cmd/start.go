package cmd

import (
	"github.com/go-macaron/bindata"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/urfave/cli/v2"
	macaron "gopkg.in/macaron.v1"

	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
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
	"gitlab.com/group-nacdlow/nacdlow-server/public"
	"gitlab.com/group-nacdlow/nacdlow-server/routes"
	routes_sim "gitlab.com/group-nacdlow/nacdlow-server/routes/simulator"
	"gitlab.com/group-nacdlow/nacdlow-server/templates"
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
		&cli.BoolFlag{
			Name:  "dev",
			Value: false,
			Usage: "enables development mode (for templates)",
		},
	},
	Action: start,
}

func getMacaron(dev bool) *macaron.Macaron {
	m := macaron.Classic()
	renderOpts := macaron.RenderOptions{
		PrefixXML:  []byte(`<?xml version="1.0" encoding="utf-8" ?>`),
		IndentJSON: true,
		IndentXML:  true,
		Funcs: []template.FuncMap{map[string]interface{}{
			"CalcSince": func(unix int64) string {
				dur := time.Since(time.Unix(unix, 0)).Round(time.Minute)
				hour := dur / time.Hour
				nomin := dur - hour*time.Hour
				min := nomin / time.Minute
				var str string
				if hour == 0 {
					str = fmt.Sprintf("%dmins", min)
				} else {
					str = fmt.Sprintf("%dhrs %dmins", hour, min)
				}
				return str
			},
			"HourStamp": func(unix int64) string {
				time := time.Unix(unix, 0)
				return time.Format("3:04pm")
			},
			"RoundPower": func(power float64) float64 {
				return math.Round(power*100) / 100
			},
		}},
	}
	staticOpts := macaron.StaticOptions{
		Expires: func() string {
			return time.Now().Add(24 * 60 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
		},
	}

	if !dev {
		renderOpts.TemplateFileSystem = bindata.Templates(bindata.Options{
			Asset:      templates.Asset,
			AssetDir:   templates.AssetDir,
			AssetNames: templates.AssetNames,
			Prefix:     "",
		})
		staticOpts.FileSystem = bindata.Static(bindata.Options{
			Asset:      public.Asset,
			AssetDir:   public.AssetDir,
			AssetNames: public.AssetNames,
			Prefix:     "",
		})
	}

	m.Use(macaron.Renderer(renderOpts))
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())
	m.Use(macaron.Static("public", staticOpts))

	m.Use(routes.ContextInit())

	m.NotFound(routes.NotFoundHandler)
	m.Get("/", routes.LoginHandler)
	m.Get("/login", routes.LoginHandler)
	m.Post("/", binding.Bind(forms.SignInForm{}), routes.PostLoginHandler)
	m.Get("/register", routes.RegisterHandler)
	m.Post("/register", binding.Bind(forms.RegisterForm{}), routes.AddUserHandler) //registers a user
	m.Get("/forgot", routes.ForgotHandler)

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
		m.Get("/internal_accounts", routes.InternalAccountsHandler)
		m.Get("/toggle_fave/:id", routes.FaveHandler)       //set device as fave
		m.Get("/remove_device/:id", routes.RemoveHandler)   //remove a device
		m.Get("/remove_room/:id", routes.RemoveRoomHandler) //removes a room

		m.Get("/move_device/:did/:rid", routes.MoveDeviceHandler) //moves device between rooms

		m.Group("/add", func() {
			m.Get("", routes.AddHandler)
			m.Get("/search_device", routes.SearchDeviceHandler)
			m.Get("/search_device/list", routes.SearchDeviceListHandler)
			m.Get("/device", routes.AddDeviceHandler)
			m.Get("/device/:id", routes.AddDeviceHandler)
			m.Post("/device", binding.Bind(forms.AddDeviceForm{}),
				routes.AddDeviceRoomPostHandler)
			m.Post("/device/:id", binding.Bind(forms.AddDeviceForm{}),
				routes.AddDeviceRoomPostHandler)
			m.Get("/device/connect/:plugin/:id", routes.ConnectDeviceHandler)
			m.Get("/device/identify/:plugin/:id", routes.IdentifyDeviceHandler)
			m.Post("/device/connect/:plugin/:id", binding.Bind(forms.AddDeviceForm{}),
				routes.ConnectDevicePostHandler)
			m.Get("/schedule", routes.AddScheduleHandler)
		})
		m.Get("/change_device/:id/:newName", routes.ChangeDeviceNameHandler) //changes the name of a device
		m.Group("/settings", func() {
			m.Get("", routes.SettingsHandler)
		})

		m.Group("", func() {
			m.Group("/add", func() {
				m.Get("/room", routes.AddRoomHandler)
				m.Post("/room", binding.Bind(forms.AddRoomForm{}),
					routes.PostRoomHandler)

			})
			m.Get("/restrict_room/:id", routes.RestrictHandler)          //restricts a room
			m.Get("/change_name/:id/:newName", routes.ChangeNameHandler) //changes the name of a room
			m.Group("/settings", func() {
				m.Group("/plugins", func() {
					m.Get("", routes.PluginsSettingsHandler)
					m.Get("/:id", routes.InstallPluginSettingsHandler)
					m.Get("/confirm/:id", routes.InstallPluginConfirmSettingsHandler) // TODO use POST so it is secure
				})

				m.Get("/plugin/:id", routes.SpecificPluginSettingsHandler)
				m.Post("/plugin/:id", routes.SpecificPluginSettingsPostHandler)

				m.Get("/plugin/:id/delete", routes.DeletePluginHandler)
				m.Get("/plugin/:id/reload", routes.ReloadPluginHandler)
				m.Get("/plugin/:id/:state([0-1])", routes.PluginStateHandler)

				m.Group("/accounts", func() {
					m.Get("", routes.AccountSettingsHandler)
					m.Post("", routes.PostAccountSettingsHandler)
					m.Get("/delete/:username", routes.DeleteAccountHandler)
					m.Post("/delete/:username", routes.PostDeleteAccountHandler)
					m.Get("/edit/:username", routes.EditAccountHandler)
					m.Post("/edit/:username", binding.Bind(forms.EditAccountForm{}),
						routes.PostEditAccountHandler)
				})

				m.Get("/devices", routes.DeviceSettingsHandler)
				m.Get("/notifications", routes.NotifictionsSettingsHandler)
				m.Get("/privacy", routes.PrivacySettingsHandler)

				m.Get("/appearance", routes.AppearanceSettingsHandler)
				m.Group("/data", func() {
					m.Get("", routes.DataSettingsHandler)
					m.Get("/json", routes.JSONDataSettingsHandler)
					m.Get("/xml", routes.XMLDataSettingsHandler)
				})
				m.Get("/appearance/font/:size", routes.FontSliderHandler)
				m.Get("/about", routes.AboutSettingsHandler)
			})
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
		m.Get("/reload_db", routes_sim.ReloadDBHandler)
		m.Get("/purge_stats", routes_sim.PurgeStatsHandler)
		m.Get("/set_main_door/:status", routes_sim.SetMainDoorHandler)
		m.Get("/set_window_status/:room/:open_count", routes_sim.SetWindowStatusHandler)
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
	m := getMacaron(clx.Bool("dev"))

	log.Printf("Starting TLS web server on :%s\n", clx.String("port"))
	server := &http.Server{Addr: fmt.Sprintf(":%s", clx.String("port")), Handler: m}
	go func() {
		log.Fatal(server.ListenAndServeTLS("fullchain.pem", "privkey.pem"))
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}()
	defer plugin.UnloadAllPlugins()

	// Capture system interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	return nil
}
