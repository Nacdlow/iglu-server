package plugin

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/Nacdlow/plugin-sdk"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	macaron "gopkg.in/macaron.v1"
)

// IgluPlugin represents a loaded Iglu plugin.
type IgluPlugin struct {
	ID            string // Unique ID, for updates, etc
	Name          string
	Author        string
	Version       string
	SettingsRoute string
	Plugin        *plugin.Plugin
	Middleware    (func() macaron.Handler)
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "IGLU_PLUGIN",
	MagicCookieValue: "MzlK0OGpIRs",
}

var pluginMap = map[string]plugin.Plugin{
	"iglu_plugin": &sdk.IgluPlugin{},
}

// LoadedPlugins is an array of all loaded plugins.
var LoadedPlugins []sdk.Iglu

// LoadPlugins will load all plugins in the `./plugins` folder.
func LoadPlugins() {

	log.Println("Loading plugins...")
	files, err := ioutil.ReadDir("./plugins")
	if err == nil {
		for _, f := range files {
			// Create an hclog.Logger
			logger := hclog.New(&hclog.LoggerOptions{
				Name:   "plugin",
				Output: os.Stdout,
				Level:  hclog.Debug,
			})
			// We're a host! Start by launching the plugin process.
			client := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: handshakeConfig,
				Plugins:         pluginMap,
				Cmd:             exec.Command(fmt.Sprintf("./plugins/%s", f.Name())),
				Logger:          logger,
			})
			defer client.Kill()

			// Connect via RPC
			rpcClient, err := client.Client()
			if err != nil {
				log.Fatal(err)
			}

			// Request the plugin
			raw, err := rpcClient.Dispense("iglu_plugin")
			if err != nil {
				log.Fatal(err)
			}

			plugin := raw.(sdk.Iglu)
			plugin.OnLoad()

			LoadedPlugins = append(LoadedPlugins, plugin)
		}
	}
	log.Printf("%d plugins loaded!\n", len(LoadedPlugins))
}
