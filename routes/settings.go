package routes

import (
	"fmt"
	"net/http"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/models/forms"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/plugin"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/tokens"

	"github.com/go-macaron/session"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

// AccountSettingsHandler handles the settings
func AccountSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Account Settings"
	var err error
	ctx.Data["Accounts"], err = models.GetUsers()
	if err != nil {
		panic(err)
	}
	ctx.Data["ArrowBack"] = 1
	ctx.HTML(http.StatusOK, "settings/accounts")
}

// PostAccountSettingsHandler handles the settings
func PostAccountSettingsHandler(ctx *macaron.Context, f *session.Flash) {
	switch ctx.Query("action") {
	case "get_invite":
		code := tokens.GenerateInviteKey()
		f.Info(fmt.Sprintf("Your new invitation code is: %s", code))
	default:
		f.Error("Unknown action")
	}
	ctx.Redirect("/settings/accounts")
}

// DeleteAccountHandler handles deleting accounts.
func DeleteAccountHandler(ctx *macaron.Context) {
	users, err := models.GetUsers()
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		if u.Username == ctx.Params("username") {
			ctx.Data["DelUser"] = u
			ctx.HTML(http.StatusOK, "settings/del_account")
			return
		}
	}
	ctx.Redirect("/settings/accounts")
}

// PostDeleteAccountHandler handles post for deleting accounts.
func PostDeleteAccountHandler(ctx *macaron.Context, f *session.Flash) {
	users, err := models.GetUsers()
	if err != nil {
		panic(err)
	}
	for _, u := range users {
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

// EditAccountHandler handles editing accounts.
func EditAccountHandler(ctx *macaron.Context) {
	users, err := models.GetUsers()
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		if u.Username == ctx.Params("username") {
			ctx.Data["EditUser"] = u
			ctx.HTML(http.StatusOK, "settings/edit_account")
			return
		}
	}
	ctx.Redirect("/settings/accounts")
}

// PostEditAccountHandler handles post for editing accounts.
func PostEditAccountHandler(ctx *macaron.Context, f *session.Flash,
	form forms.EditAccountForm) {
	users, err := models.GetUsers()
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		if u.Username == form.Email {
			var updateCols []string
			user := models.User{Username: form.Email}

			// Update password
			if form.Password != "" {
				if form.Password != form.RePassword {
					f.Error("Passwords do not match!")
					ctx.Redirect(fmt.Sprintf("/settings/accounts/edit/%s", form.Email))
					return
				}
				pass, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
				if err != nil {
					panic(err)
				}

				updateCols = append(updateCols, "password")
				user.Password = string(pass)
			}

			// Update first name
			if form.FirstName != "" && form.FirstName != u.FirstName {
				updateCols = append(updateCols, "first_name")
				user.FirstName = form.FirstName
			}

			// Update first name
			if form.LastName != "" && form.LastName != u.LastName {
				updateCols = append(updateCols, "last_name")
				user.LastName = form.LastName
			}

			err := models.UpdateUserCols(&user, updateCols...)
			if err != nil {
				f.Error("Failed to update user!")
			} else {
				f.Success("User updated!")
			}
		}
	}
	ctx.Redirect("/settings/accounts")
}

// AppearanceSettingsHandler handles the settings
func AppearanceSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Appearance Settings"
	ctx.Data["ArrowBack"] = 1
	ctx.HTML(http.StatusOK, "settings/appearance")
}

// PluginsSettingsHandler handles the settings
func PluginsSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Installed Plugins"
	ctx.Data["Plugins"] = plugin.LoadedPlugins
	ctx.Data["ArrowBack"] = 1
	ctx.HTML(http.StatusOK, "settings/plugins")
}

// SettingsHandler handles the settings
func SettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Settings"
	ctx.HTML(http.StatusOK, "settings")
}

// LibraryDesc represents the description of an open-source library.
type LibraryDesc struct {
	Author     string
	ProjectURL string
}

// Libraries contains the list of the libraries used in this application.
var Libraries = map[string]LibraryDesc{
	"MDBootstrap":        {ProjectURL: "https://mdbootstrap.com/"},
	"jQuery":             {ProjectURL: "https://jquery.com/"},
	"Skycons":            {Author: "Dark Sky", ProjectURL: "https://darkskyapp.github.io/skycons/"},
	"Polyfill.io":        {Author: "Financial Times", ProjectURL: "https://polyfill.io"},
	"Popper.js":          {Author: "Federico Zivolo", ProjectURL: "https://popper.js.org/"},
	"Macaron":            {Author: "Jiahua Chen (Unknwon)", ProjectURL: "https://github.com/go-macaron/macaron"},
	"xorm":               {ProjectURL: "https://gitea.com/xorm/xorm"},
	"cli":                {Author: "Jeremy Saenz", ProjectURL: "https://github.com/urfave/cli"},
	"viper":              {Author: "Steve Francia", ProjectURL: "https://github.com/spf13/viper"},
	"Go Dark Sky API":    {Author: "Aaron Longwell", ProjectURL: "https://github.com/adlio/darksky"},
	"TOML parser for Go": {Author: "Andrew Gallant", ProjectURL: "https://github.com/BurntSushi/toml"},
	"gofakeit":           {Author: "Brian Voelker", ProjectURL: "https://github.com/brianvoe/gofakeit"},
}

// AboutSettingsHandler handles the about settings page
func AboutSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "About"
	ctx.Data["Libraries"] = Libraries
	ctx.Data["ArrowBack"] = 1
	ctx.HTML(http.StatusOK, "settings/about")
}
