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

func DeleteAccountHandler(ctx *macaron.Context) {
	for _, u := range models.GetUsers() {
		if u.Username == ctx.Params("username") {
			ctx.Data["DelUser"] = u
			ctx.HTML(200, "settings/del_account")
			return
		}
	}
	ctx.Redirect("/settings/accounts")
}

func PostDeleteAccountHandler(ctx *macaron.Context, f *session.Flash) {
	for _, u := range models.GetUsers() {
		if u.Username == ctx.Query("username") {
			err := models.DeleteUser(ctx.Query("username"))
			if err != nil {
				f.Error("Failed to delete user!")
			} else {
				f.Success("User deleted!")
			}
		}
	}
	ctx.Redirect("/settings/accounts")
}

func EditAccountHandler(ctx *macaron.Context) {
	for _, u := range models.GetUsers() {
		if u.Username == ctx.Params("username") {
			ctx.Data["EditUser"] = u
			ctx.HTML(200, "settings/edit_account")
			return
		}
	}
	ctx.Redirect("/settings/accounts")
}

func PostEditAccountHandler(ctx *macaron.Context, f *session.Flash) {
	for _, u := range models.GetUsers() {
		if u.Username == ctx.Query("username") {
			//models.UpdateUser(...
		}
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
