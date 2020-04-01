package routes

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/Nacdlow/plugin-sdk"
	"github.com/go-macaron/session"
	"github.com/hashicorp/go-getter"
	macaron "gopkg.in/macaron.v1"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
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
	appendExe := ""
	if runtime.GOOS == "windows" {
		appendExe = ".exe"
	}
	repo := fmt.Sprintf("%s/%s/%s%s.xz?checksum=file:%s/%s/%s%s.xz.sha256sum&archive=xz",
		repoURL,
		platform,
		id,
		appendExe,
		repoURL,
		platform,
		id,
		appendExe)
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

	if runtime.GOOS == "windows" {
		plugin.LoadPlugin(id + ".exe")
	} else {
		plugin.LoadPlugin(id)
	}

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

type FilledPluginConfig struct {
	Conf  sdk.PluginConfig
	Value string
}

func SpecificPluginSettingsHandler(ctx *macaron.Context, f *session.Flash) {
	var pl *plugin.IgluPlugin
	plugins := plugin.GetLoadedPlugins()
	for i := range plugins {
		if plugins[i].Manifest.Id == ctx.Params("id") {
			pl = &plugins[i]
			break
		}
	}

	if pl == nil {
		f.Error("Cannot find plugin.")
		ctx.Redirect("/settings/plugins")
		return
	}

	user := ctx.Data["User"].(*models.User)

	var filledConfigs []FilledPluginConfig

	for i, c := range pl.Config {
		if !c.IsUserSpecific && ctx.Data["IsAdmin"].(int) != 1 {
			continue
		}
		kvName := fmt.Sprintf("%s.%s", pl.Manifest.Id, c.Key)

		filled := FilledPluginConfig{Conf: pl.Config[i]}

		if c.IsUserSpecific {
			val, ok := user.PluginKVStore[kvName]
			if ok {
				filled.Value = val
			}
		} else {
			// TODO add global load
		}
		filledConfigs = append(filledConfigs, filled)
	}

	ctx.Data["Plugin"] = pl
	ctx.Data["FilledConfigs"] = filledConfigs
	ctx.HTML(http.StatusOK, "settings/plugin_setting")
}

func SpecificPluginSettingsPostHandler(ctx *macaron.Context, f *session.Flash) {
	var pl *plugin.IgluPlugin
	plugins := plugin.GetLoadedPlugins()
	for i := range plugins {
		if plugins[i].Manifest.Id == ctx.Params("id") {
			pl = &plugins[i]
			break
		}
	}

	if pl == nil {
		f.Error("Cannot find plugin.")
		ctx.Redirect("/settings/plugins")
		return
	}

	user := ctx.Data["User"].(*models.User)
	userUpdated := false
	if user.PluginKVStore == nil {
		user.PluginKVStore = make(map[string]string)
	}

	for _, c := range pl.Config {
		if !c.IsUserSpecific && ctx.Data["IsAdmin"].(int) != 1 {
			continue
		}

		fieldName := fmt.Sprintf("field-%s", c.Key)
		var input string

		switch c.Type {
		case sdk.OptionValue:
			fallthrough
		case sdk.StringValue:
			input = ctx.Query(fieldName)
		case sdk.BooleanValue:
			input = fmt.Sprintf("%t", ctx.QueryBool(fieldName))
		case sdk.NumberValue:
			input = strconv.Itoa(ctx.QueryInt(fieldName))
		}
		if len(input) > 0 {
			if c.IsUserSpecific {
				kvName := fmt.Sprintf("%s.%s", pl.Manifest.Id, c.Key)
				user.PluginKVStore[kvName] = input

				userUpdated = true
			} else {
				// TODO add global save
			}
		}
	}
	if userUpdated {
		models.UpdateUserCols(user, "plugin_kv_store")
	}
	// TODO send the KVs to the plugin
}

func DeletePluginHandler(ctx *macaron.Context, f *session.Flash) {
	if !plugin.IsPluginLoaded(ctx.Params("id")) {
		f.Error("Cannot find plugin.")
		ctx.Redirect("/settings/plugins")
		return
	}

	plugin.DeletePlugin(ctx.Params("id"))
	ctx.Redirect("/settings/plugins")
}

func ReloadPluginHandler(ctx *macaron.Context, f *session.Flash) {
	if !plugin.IsPluginLoaded(ctx.Params("id")) {
		f.Error("Cannot find plugin.")
		ctx.Redirect("/settings/plugins")
		return
	}

	plugin.ReloadPlugin(ctx.Params("id"))
	f.Success("Plugin reloaded!")
	ctx.Redirect("/settings/plugin/" + ctx.Params("id"))
}

func PluginStateHandler(ctx *macaron.Context) {
	id := ctx.Params("id")
	_, err := plugin.GetPlugin(id)
	if err != nil {
		return
	}

	switch ctx.ParamsInt("state") {
	case 0:
		plugin.StopPlugin(id)
	case 1:
		plugin.StartPlugin(id)
	}
	ctx.Redirect("/settings/plugins")
}
