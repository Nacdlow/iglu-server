package routes

import (
	"fmt"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/plugin"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/tokens"

	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

// AccountSettingsHandler handles the settings
func AccountSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Account Settings"
	ctx.Data["Accounts"] = models.GetUsers()
	ctx.HTML(200, "settings/accounts")
}

// PostAccountSettingsHandler handles the settings
func PostAccountSettingsHandler(ctx *macaron.Context, f *session.Flash) {
	switch ctx.Query("action") {
	case "get_invite":
		code := tokens.GenerateInviteKey()
		f.Info(fmt.Sprintf("Your new invitation code is: %s", code))
		break
	default:
		f.Error("Unknown action")
	}
	ctx.Redirect("/settings/accounts")
}

// AppearanceSettingsHandler handles the settings
func AppearanceSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Appearance Settings"
	ctx.HTML(200, "settings/appearance")
}

// PluginsSettingsHandler handles the settings
func PluginsSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Plugins"
	ctx.Data["Plugins"] = plugin.LoadedPlugins
	ctx.HTML(200, "settings/plugins")
}

// SettingsHandler handles the settings
func SettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Settings"
	ctx.HTML(200, "settings")
}
