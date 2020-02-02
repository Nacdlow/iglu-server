package plugin

import (
	"fmt"
	macaron "gopkg.in/macaron.v1"
	"io/ioutil"
	"log"
	"plugin"
	"runtime"
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

// LoadedPlugins is an array of all loaded plugins.
var LoadedPlugins []IgluPlugin

// LoadPlugins will load all plugins in the `./plugins` folder.
func LoadPlugins() {
	if !(runtime.GOOS == "linux" || runtime.GOOS == "darwin") {
		log.Println("Plugins only supported on Linux and macOS!")
		log.Println("Plugins will not be loaded.")
		return
	}
	log.Println("Loading plugins...")
	files, err := ioutil.ReadDir("./plugins")
	if err == nil {
		for _, f := range files {
			p, err := plugin.Open(fmt.Sprintf("./plugins/%s", f.Name()))
			if err != nil {
				log.Printf("Failed to load plugin '%s'!\n", f.Name())
				log.Println(err)
				continue
			}
			id, err := p.Lookup("ID")
			if err != nil {
				log.Printf("Failed to load variables for plugin '%s'!\n", f.Name())
				continue
			}
			name, err := p.Lookup("NAME")
			if err != nil {
				log.Printf("Failed to load variables for plugin '%s'!\n", f.Name())
				continue
			}
			author, err := p.Lookup("AUTHOR")
			if err != nil {
				log.Printf("Failed to load variables for plugin '%s'!\n", f.Name())
				continue
			}
			version, err := p.Lookup("VERSION")
			log.Printf("Loading %s %s...", *name.(*string), *version.(*string))
			load, err := p.Lookup("Load")
			if err != nil {
				log.Printf("Failed to run Load for plugin '%s'!\n", f.Name())
				continue
			}
			load.(func())()
			pl := IgluPlugin{
				ID:      *id.(*string),
				Name:    *name.(*string),
				Author:  *author.(*string),
				Version: *version.(*string),
			}
			route, err := p.Lookup("ROUTE")
			if err == nil {
				pl.SettingsRoute = *route.(*string)
			}
			middleware, err := p.Lookup("Middleware")
			if err == nil {
				pl.Middleware = middleware.(func() macaron.Handler)
			}
			LoadedPlugins = append(LoadedPlugins, pl)
		}
	}
	log.Printf("%d plugins loaded!\n", len(LoadedPlugins))
}
