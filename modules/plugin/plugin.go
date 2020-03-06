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

type PluginState int64

const (
	Stopped = iota
	Running
	Crashed
)

// IgluPlugin represents a loaded Iglu plugin.
type IgluPlugin struct {
	Plugin   sdk.Iglu
	client   *plugin.Client
	State    PluginState
	Filename string
	Config   []sdk.PluginConfig
	Manifest sdk.PluginManifest
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "IGLU_PLUGIN",
	MagicCookieValue: "MzlK0OGpIRs",
}

var pluginMap = map[string]plugin.Plugin{
	"iglu_plugin": &sdk.IgluPlugin{},
}

var loadedPlugins []IgluPlugin
var mutex sync.Mutex

func markCrashedPlugins() {
	mutex.Lock()
	defer mutex.Unlock()
	for i, plugin := range loadedPlugins {
		if plugin.client.Exited() && plugin.State == Running {
			log.Printf("Plugin '%s' crashed\n", plugin.Manifest.Id)
			loadedPlugins[i].State = Crashed
		}
	}
}

// GetLoadedPlugins returns a list of loaded and running plugins.
func GetLoadedPlugins() []IgluPlugin {
	markCrashedPlugins()
	return loadedPlugins
}

// DeletePlugin will unload and delete plugin from disk.
func DeletePlugin(id string) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, plugin := range loadedPlugins {
		if plugin.client != nil && plugin.client.Exited() {
			plugin.client.Kill()
		}
		os.Remove("./plugins/" + plugin.Filename)
		loadedPlugins = append(loadedPlugins[:i], loadedPlugins[i+1:]...)
		return
	}
}

// UnloadAllPlugins unloads all loaded plugins.
func UnloadAllPlugins() {
	mutex.Lock()
	defer mutex.Unlock()
	log.Println("Unloading plugins...")
	for i, plugin := range loadedPlugins {
		if plugin.client != nil && !plugin.client.Exited() {
			plugin.client.Kill()
		}
		loadedPlugins[i].State = Stopped
	}
	log.Println("Plugins unloaded")
}

// UnloadPlugin unloads a loaded plugin.
func UnloadPlugin(id string) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, plugin := range loadedPlugins {
		if plugin.client != nil {
			plugin.client.Kill()
		}
		loadedPlugins[i].State = Stopped
		return
	}
}

// GetPlugin returns a loaded plugin.
func GetPlugin(id string) (*IgluPlugin, error) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, plugin := range loadedPlugins {
		if plugin.client != nil {
			return &loadedPlugins[i], nil
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
		Managed:         true,
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

	loadedPlugins = append(loadedPlugins, IgluPlugin{
		Plugin:   plugin,
		client:   client,
		State:    Running,
		Filename: f,
		Manifest: plugin.GetManifest(),
		Config:   plugin.GetPluginConfiguration(),
	})
}

// IsPluginLoaded returns whether a plugin is loaded based on the plugin ID.
// A positive doesn't mean the plugin is running.
func IsPluginLoaded(id string) bool {
	for _, pl := range loadedPlugins {
		if pl.Plugin != nil && pl.Manifest.Id == id {
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
	log.Printf("%d plugins loaded!\n", len(loadedPlugins))
}
