package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/Nacdlow/plugin-sdk"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// IgluPlugin represents a loaded Iglu plugin.
type IgluPlugin struct {
	Plugin sdk.Iglu
	client *plugin.Client
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
var LoadedPlugins []IgluPlugin
var mutex sync.Mutex

// UnloadAllPlugins unloads all loaded plugins.
func UnloadAllPlugins() {
	mutex.Lock()
	defer mutex.Unlock()
	log.Println("Unloading plugins...")
	for _, plugin := range LoadedPlugins {
		if plugin.client != nil {
			plugin.client.Kill()
		}
	}
	LoadedPlugins = []IgluPlugin{}
	log.Println("Plugins unloaded")
}

// UnloadPlugin unloads a loaded plugin.
func UnloadPlugin(id string) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, plugin := range LoadedPlugins {
		if plugin.client != nil {
			plugin.client.Kill()
		}
		LoadedPlugins = append(LoadedPlugins[:i], LoadedPlugins[i+1:]...)
		return
	}
}

// GetPlugin returns a loaded plugin.
func GetPlugin(id string) (*IgluPlugin, error) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, plugin := range LoadedPlugins {
		if plugin.client != nil {
			return &LoadedPlugins[i], nil
		}
	}
	return nil, errors.New("Plugin is not loaded")
}

// LoadPlugin loads a plugin from a file path.
func LoadPlugin(f string) {
	mutex.Lock()
	defer mutex.Unlock()
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
		Cmd:             exec.Command(fmt.Sprintf("./plugins/%s", f)),
		Logger:          logger,
	})

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
	err = plugin.OnLoad()
	if err != nil {
		log.Printf("Failed to load plugin %s (onLoad): %s\n", f, err)
	}

	// We cannot load the same plugin twice!
	if IsPluginLoaded(plugin.GetManifest().Id) {
		log.Printf("Cannot load plugin '%s' as it is already loaded!\n", plugin.GetManifest().Id)
		client.Kill()
		return
	}

	LoadedPlugins = append(LoadedPlugins, IgluPlugin{
		Plugin: plugin,
		client: client,
	})
}

// IsPluginLoaded returns whether a plugin is loaded based on the plugin ID.
func IsPluginLoaded(id string) bool {
	for _, pl := range LoadedPlugins {
		if pl.Plugin != nil && pl.Plugin.GetManifest().Id == id {
			return true
		}
	}
	return false
}

// LoadPlugins will load all plugins in the `./plugins` folder.
func LoadPlugins() {
	log.Println("Loading plugins...")
	files, err := ioutil.ReadDir("./plugins")
	if err == nil {
		for _, f := range files {
			LoadPlugin(f.Name())
		}
	}
	log.Printf("%d plugins loaded!\n", len(LoadedPlugins))
}
