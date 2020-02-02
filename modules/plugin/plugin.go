package plugin

import (
	"fmt"
	"io/ioutil"
	"log"
	"plugin"
)

var plugins map[string]*plugin.Plugin

func LoadPlugins() {
	log.Println("Loading plugins...")
	files, err := ioutil.ReadDir("./plugins")
	var pluginCount int
	if err == nil {
		for _, f := range files {
			p, err := plugin.Open(fmt.Sprintf("./plugins/%s", f.Name()))
			if err != nil {
				log.Printf("Failed to load plugin '%s'!\n", f.Name())
				continue
			}
			load, err := p.Lookup("Load")
			if err != nil {
				log.Printf("Failed to run Load for plugin %s!\n", f.Name())
				continue
			}
			load.(func())()
			pluginCount++
		}
	}
	log.Printf("%d plugins loaded!\n", pluginCount)

}
