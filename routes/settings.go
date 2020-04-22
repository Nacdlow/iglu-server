package routes

import (
	"fmt"
	"net/http"

	"github.com/Nacdlow/iglu-server/models"
	"github.com/Nacdlow/iglu-server/models/forms"
	"github.com/Nacdlow/iglu-server/modules/plugin"
	"github.com/Nacdlow/iglu-server/modules/tokens"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

// AccountSettingsHandler handles the settings
func AccountSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Accounts"
	var err error
	ctx.Data["Accounts"], err = models.GetUsers()
	if err != nil {
		panic(err)
	}
	ctx.Data["ArrowBack"] = 1
	ctx.HTML(http.StatusOK, "settings/accounts")
}

func SetAvatarHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	users, err := models.GetUsers()
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		if u.Username == ctx.Params("username") {
			u.Avatar = "/img/profiles/" + ctx.Params("avatar") + ".jpg"

			err = models.UpdateUserCols(&u, "Avatar")
			if err != nil {
				panic(err)
			}
			f.Success("Avatar updated")
			ctx.Redirect("/settings/accounts")
			return
		}
	}

	f.Error("User does not exist.")
	ctx.Redirect("/settings/accounts")
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
	form forms.EditAccountForm, errs binding.Errors) {
	if len(errs) > 0 {
		f.Error("Missing required fields!")
		ctx.Redirect("/settings/accounts")
		return
	}
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

			// Update Role
			updateCols = append(updateCols, "role")
			user.Role = models.UserRole(form.Role)

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

// EditAvatarHandler handles editing account avatars.
func EditAvatarHandler(ctx *macaron.Context) {
	ctx.Data["ArrowBack"] = 1
	ctx.Data["uname"] = ctx.Params("username")
	ctx.HTML(http.StatusOK, "settings/edit_avatar")
}

// AppearanceSettingsHandler handles the settings
func AppearanceSettingsHandler(ctx *macaron.Context, sess session.Store) {
	ctx.Data["NavTitle"] = "Appearance"
	ctx.Data["ArrowBack"] = 1

	if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
		switch user.FontSize {
		case "xsmall":
			ctx.Data["SliderNumber"] = 1
			ctx.Data["SliderText"] = "Tiny"
		case "small":
			ctx.Data["SliderNumber"] = 2
			ctx.Data["SliderText"] = "Small"
		case "medium":
			ctx.Data["SliderNumber"] = 3
			ctx.Data["SliderText"] = "Normal"
		case "large":
			ctx.Data["SliderNumber"] = 4
			ctx.Data["SliderText"] = "Large"
		case "xlarge":
			ctx.Data["SliderNumber"] = 5
			ctx.Data["SliderText"] = "Huge"
		default:
			ctx.Data["SliderNumber"] = 3
			ctx.Data["SliderText"] = "Normal"
		}

		ctx.HTML(http.StatusOK, "settings/appearance")
	}
}

// PluginsSettingsHandler handles the settings
func PluginsSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Installed Plugins"
	ctx.Data["Plugins"] = plugin.GetLoadedPlugins()
	ctx.Data["ArrowBack"] = 1
	ctx.HTML(http.StatusOK, "settings/plugins")
}

// SettingsHandler handles the settings
func SettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Settings"
	ctx.HTML(http.StatusOK, "settings")
}

// DeviceSettingsHandler handles the device settings
func DeviceSettingsHandler(ctx *macaron.Context) {
	var err error
	ctx.Data["Devices"], err = models.GetDevices()
	if err != nil {
		panic(err)
	}
	ctx.Data["ArrowBack"] = 1

	ctx.Data["NavTitle"] = "Device Settings"
	ctx.HTML(http.StatusOK, "settings/devices")
}

// NotificationsSettingsHandler handles the device settings
func NotifictionsSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Device Settings"
	ctx.HTML(http.StatusOK, "settings/notifications")
}

// PrivacySettingsHandler handles the device settings
func PrivacySettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Device Settings"
	ctx.HTML(http.StatusOK, "settings/privacy")
}

// HelpHandler handles the how-to help section.
func HelpHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Device Settings"
	ctx.HTML(http.StatusOK, "settings/help")
}

// HelpHandler handles the how-to help section/fave device.
func HelpFdHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "How to Favourite a Device"
	ctx.HTML(http.StatusOK, "settings/how-to/fave_device")
}

// HelpHandler handles the how-to help section./add device
func HelpAdHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "How to Add a Device"
	ctx.HTML(http.StatusOK, "settings/how-to/add_device")
}

// HelpHandler handles the how-to help section./add room
func HelpArHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "How to Add a Room"
	ctx.HTML(http.StatusOK, "settings/how-to/add_room")
}

// HelpHandler handles the how-to help section./add user
func HelpAuHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "How to Add a User"
	ctx.HTML(http.StatusOK, "settings/how-to/add_user")
}

// HelpHandler handles the how-to help section./delete device
func HelpDdHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "How to Delete a Device"
	ctx.HTML(http.StatusOK, "settings/how-to/delete_device")
}

// HelpHandler handles the how-to help section./delete room
func HelpDrHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "How to Delete a Room"
	ctx.HTML(http.StatusOK, "settings/how-to/delete_room")
}

// HelpHandler handles the how-to help section./edit room
func HelpErHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "How to Edit Room Name"
	ctx.HTML(http.StatusOK, "settings/how-to/edit_room")
}

// HelpHandler handles the how-to help section./restrict rooms
func HelpRrHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Restrict/Un-restrict a Room"
	ctx.HTML(http.StatusOK, "settings/how-to/restrict_room")
}

// HelpHandler handles the how-to help section./change password
func HelpCpHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Change Password"
	ctx.HTML(http.StatusOK, "settings/how-to/change_password")
}

// HelpHandler handles the how-to help section./remove user
func HelpRuHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Remove a user"
	ctx.HTML(http.StatusOK, "settings/how-to/remove_user")
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

func DataSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Import/Export"
	ctx.Data["ArrowBack"] = 1
	ctx.HTML(http.StatusOK, "settings/data")
}

func JSONDataSettingsHandler(ctx *macaron.Context) {
	alerts, _ := models.GetAlerts()
	devices, _ := models.GetDevices()
	rooms, _ := models.GetRooms()
	room_stats, _ := models.GetRoomStats()
	schedules, _ := models.GetSchedules()
	statistics, _ := models.GetStats()
	ctx.JSON(200, map[string]interface{}{
		"alerts":     alerts,
		"device":     devices,
		"rooms":      rooms,
		"roomStats":  room_stats,
		"schedules":  schedules,
		"statistics": statistics,
	})
}
func XMLDataSettingsHandler(ctx *macaron.Context) {
	alerts, _ := models.GetAlerts()
	devices, _ := models.GetDevices()
	rooms, _ := models.GetRooms()
	room_stats, _ := models.GetRoomStats()
	schedules, _ := models.GetSchedules()
	statistics, _ := models.GetStats()
	type Export struct {
		Alerts     []models.Alert  `xml:"alerts"`
		Devices    []models.Device `xml:"devices"`
		Rooms      []models.Room
		Room_stats []models.RoomStat
		Schedules  []models.Schedule
		Statistics []models.Statistic
	}
	export := Export{alerts, devices, rooms, room_stats, schedules, statistics}
	ctx.XML(200, export)
}
