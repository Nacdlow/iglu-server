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
	"github.com/hashicorp/go-getter"
	macaron "gopkg.in/macaron.v1"

	"gitlab.com/group-nacdlow/nacdlow-server/modules/plugin"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
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
	repoURL := settings.Config.Get("Marketplace.RepositoryURL")
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	repo := fmt.Sprintf("%s/%s/%s.xz?checksum=file:%s/%s/%s.sha256sum&archive=xz",
		repoURL,
		platform,
		id,
		repoURL,
		platform,
		id)
	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("Failed to get wd: %s", err))
	}
	pluginFile := fmt.Sprintf("./plugins/%s", id)

	client := &getter.Client{
		Src:  repo,
		Dst:  pluginFile,
		Pwd:  pwd,
		Mode: getter.ClientModeFile,
	}

	err = client.Get()
	if err != nil {
		f.Error("Failed to fetch plugin.")
		fmt.Println(err)
		ctx.Redirect("/settings/plugins")
		return
	}
	os.Chmod(pluginFile, 0700)

	plugin.LoadPlugin(id)

	f.Success("Plugin installed and loaded!")
	ctx.Redirect("/settings/plugins")
}

// InstallPluginSettingsHandler handles installing the plugin.
func InstallPluginSettingsHandler(ctx *macaron.Context, f *session.Flash) {
	id := ctx.Params("id")
	repoDesc := fmt.Sprintf("%s/%s-%s/%s.toml",
		settings.Config.Get("Marketplace.RepositoryURL"),
		runtime.GOOS,
		runtime.GOARCH,
		id)
	desc, err := getPluginDesc(repoDesc)
	if err != nil {
		f.Error("Cannot find plugin.")
		ctx.Redirect("/settings/plugins")
		return
	}
	ctx.Data["Plugin"] = desc

	ctx.HTML(http.StatusOK, "settings/plugin_install_confirm")
}
