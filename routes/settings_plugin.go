package routes

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"

	"gitlab.com/group-nacdlow/nacdlow-server/modules/plugin"
)

type pluginDescription struct {
	ID      string `toml:"ID"`
	Name    string `toml:"NAME"`
	Author  string `toml:"AUTHOR"`
	Version string `toml:"VERSION"`
}

func getPluginDesc(url string) (pluginDescription, error) {
	resp, err := http.Get(url)
	if err != nil {
		return pluginDescription{}, err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pluginDescription{}, err
	}
	var desc pluginDescription
	if _, err = toml.Decode(string(html), &desc); err != nil {
		return pluginDescription{}, err
	}
	return desc, nil
}

func downloadPlugin(name, url string) error {
	// Download plugin binary
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Place in ./plugins folder
	file := fmt.Sprintf("./plugins/%s", name)
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// InstallPluginConfirmSettingsHandler handles installing the plugin.
func InstallPluginConfirmSettingsHandler(ctx *macaron.Context, f *session.Flash) {
	id := ctx.Params("id")
	repo := fmt.Sprintf("https://market.nacdlow.com/repo/%s-%s/%s", runtime.GOOS, runtime.GOARCH, id)
	err := downloadPlugin(id, repo)
	if err != nil {
		panic(err)
	}

	f.Success("Plugin installed! Please restart the server to load the plugin.")
	ctx.Redirect("/settings/plugins")
}

// InstallPluginSettingsHandler handles installing the plugin.
func InstallPluginSettingsHandler(ctx *macaron.Context) {
	id := ctx.Params("id")
	repoDesc := fmt.Sprintf("https://market.nacdlow.com/repo/%s-%s/%s.toml", runtime.GOOS, runtime.GOARCH, id)
	desc, err := getPluginDesc(repoDesc)
	if err != nil {
		panic(err)
	}
	ctx.Data["Plugin"] = desc

	ctx.HTML(http.StatusOK, "settings/plugin_install_confirm")
}

// PluginSettingPage handles rendering the plugin pages.
func PluginSettingPage(ctx *macaron.Context, f *session.Flash) {
	if !plugin.IsPluginLoaded(ctx.Params("id")) {
		f.Error("Plugin not loaded!")
		ctx.Redirect("/settings/plugins")
		return
	}

	pl, err := plugin.GetPlugin("test")
	if err != nil {
		panic(err)
	}
	pl.Plugin.PluginHTTP(ctx.Resp, ctx.Req.Request)
	return
}
