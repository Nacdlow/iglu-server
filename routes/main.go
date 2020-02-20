package routes

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/models/forms"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/plugin"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/simulation"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/tokens"

	"github.com/BurntSushi/toml"
	"github.com/go-macaron/session"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

// NotFoundHandler handles 404 errors
func NotFoundHandler(ctx *macaron.Context) {
	ctx.HTML(404, "notfound")
}

// DashboardHandler handles rendering the dashboard.
func DashboardHandler(ctx *macaron.Context, f *session.Flash) {
	ctx.Data["NavTitle"] = "Dashboard"
	ctx.Data["IsDashboard"] = 1
	if simulation.Env.ForecastData != nil {
		ctx.Data["Temperature"] = math.Round(simulation.Env.ForecastData.Currently.Temperature)
		ctx.Data["Summary"] = simulation.Env.ForecastData.Currently.Summary
		icon := simulation.Env.ForecastData.Currently.Icon
		icon = strings.ToUpper(icon)
		icon = strings.ReplaceAll(icon, "-", "_")
		ctx.Data["WeatherIcon"] = template.JS(icon)
	}
	ctx.Data["Devices"] = models.GetDevices()
	ctx.HTML(200, "dashboard")
}

func BatteryStatHandler(ctx *macaron.Context) {
	perc := int64((simulation.Env.BatteryStore / settings.Config.GetFloat64("Simulation.BatteryCapacityKWH")) * 100)
	if perc < 0 {
		perc = 0
	}
	if perc < 15 {
		ctx.Data["BatState"] = 0 // empty
	} else if perc < 50 {
		ctx.Data["BatState"] = 1 // quarter
	} else if perc < 75 {
		ctx.Data["BatState"] = 2 // half
	} else if perc < 95 {
		ctx.Data["BatState"] = 3 // three-quarter
	} else {
		ctx.Data["BatState"] = 4 // full
	}

	ctx.Data["BatteryPerc"] = perc
	ctx.HTML(200, "battery")
}

// SettingsHandler handles the settings
func SettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Settings"
	ctx.HTML(200, "settings")
}

// AddHandler handles the add page
func AddHandler(ctx *macaron.Context) {
	ctx.Data["CrossBack"] = 1
	ctx.HTML(200, "add")
}

// AddRoomHandler handles the add room page
func AddRoomHandler(ctx *macaron.Context) {
	ctx.Data["BackLink"] = "/add"
	ctx.Data["Rooms"] = models.GetRooms()
	ctx.HTML(200, "add_room")
}

// AlertsHandler handles the alerts page
func AlertsHandler(ctx *macaron.Context) {
	ctx.Data["CrossBack"] = 1
	ctx.HTML(200, "alerts")
}

// AddDeviceHandler handles the add device page
func AddDeviceHandler(ctx *macaron.Context) {
	ctx.Data["BackLink"] = "/add"
	ctx.Data["Rooms"] = models.GetRooms()
	if ctx.Params("id") != "" {
		ctx.Data["RoomSelected"] = 1
		ctx.Data["RoomID"] = ctx.Params("id")
	}
	ctx.HTML(200, "add_device")
}

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

	ctx.HTML(200, "settings/plugin_install_confirm")
}

// SpecificRoomsHandler handles the specific rooms
func SpecificRoomsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = ctx.Params("roomType")
	if ctx.Params("name") == "bedrooms" {
		ctx.Data["Bedrooms"] = 1
		ctx.Data["Rooms"] = models.GetRooms()
	} else if ctx.Params("name") == "bathrooms" {
		ctx.Data["Bathrooms"] = 1
		ctx.Data["Rooms"] = models.GetRooms()
	} else {
		room, err := models.GetRoom(ctx.ParamsInt64("name"))
		if err != nil {
			ctx.Redirect("/rooms")
			return
		}
		ctx.Data["Room"] = room
		ctx.Data["Devices"] = models.GetDevices()
	}

	ctx.Data["ArrowBack"] = 1
	ctx.Data["IsRooms"] = 1
	ctx.Data["IsSpecificRoom"] = 1
	ctx.HTML(200, "specificRooms")
}

// OverviewHandler handles the overview page
func OverviewHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Overview"
	ctx.Data["IsOverview"] = 1
	ctx.Data["Devices"] = models.GetDevices()
	ctx.HTML(200, "overview")
}

func AddDeviceRoomPostHandler(ctx *macaron.Context, form forms.AddDeviceForm, f *session.Flash) {
	device := &models.Device{
		RoomID:      form.RoomID,
		Type:        models.DeviceType(form.DeviceType),
		Description: form.Description,
		ToggledUnix: time.Now().Unix(),
	}
	for _, d := range models.GetDevices() {
		if d.RoomID == device.RoomID && d.Type == device.Type {
			switch device.Type {
			case models.Light: // We can't have two main lights
				if form.IsMainLight && device.IsMainLight {
					f.Error("There can only be one main light per room.")
					ctx.Redirect(fmt.Sprintf("/room/%d", device.RoomID))
					return
				}
			case models.TempControl: // We can't have two temperature controls
				f.Error("There can only be one temperature control per room.")
				ctx.Redirect(fmt.Sprintf("/room/%d", device.RoomID))
				return
			}
		}
	}

	switch device.Type {
	case models.Light:
		device.Brightness = 10
		if form.IsMainLight {
			device.IsMainLight = true
		}
	case models.TempControl:
		device.Temp = 22.0
	case models.Speaker:
		device.Volume = 8
	}
	models.AddDevice(device)
	ctx.Redirect(fmt.Sprintf("/room/%d", form.RoomID))
}

// RoomsHandler handles rendering the rooms page
func RoomsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Rooms"
	ctx.Data["IsRooms"] = 1
	rooms := models.GetRooms()
	for i, _ := range rooms {
		rooms[i].LoadMainDevices()
	}
	ctx.Data["Rooms"] = rooms
	ctx.HTML(200, "rooms")
}

// PostRoomHandler handles post request for room page, to add a room.
func PostRoomHandler(ctx *macaron.Context, form forms.AddRoomForm) {
	room := &models.Room{
		RoomName:    form.RoomName,
		Description: form.Description,
		RoomType:    models.RType(form.RoomType),
		WindowCount: form.WindowsCount,
	}
	if form.PartOfRoom >= 0 {
		room.IsSubRoom = true
		room.PartOfRoom = form.PartOfRoom
	}
	models.AddRoom(room)
	ctx.Redirect("/rooms")
}

// DevicesHandler handles the devices page
func DevicesHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Devices"
	ctx.HTML(200, "devices")
}

func ToggleHandler(ctx *macaron.Context) {
	for _, device := range models.GetDevices() {
		if device.DeviceID == ctx.ParamsInt64("id") {
			if device.Type == models.Light || device.Type == models.Speaker ||
				device.Type == models.TempControl || device.Type == models.TV {
				models.UpdateDeviceCols(&models.Device{
					DeviceID:    device.DeviceID,
					Status:      !device.Status,
					ToggledUnix: time.Now().Unix()}, "status", "toggled_unix")
			}
			break
		}
	}
}

func SliderHandler(ctx *macaron.Context) {
	for _, device := range models.GetDevices() {
		if device.DeviceID == ctx.ParamsInt64("id") {
			if device.Type == models.Speaker {
				models.UpdateDeviceCols(&models.Device{
					DeviceID: device.DeviceID,
					Volume:   ctx.ParamsInt64("value")}, "volume")
			} else if device.Type == models.TempControl {
				models.UpdateDeviceCols(&models.Device{
					DeviceID: device.DeviceID,
					Temp:     ctx.ParamsFloat64("value")}, "temp")
			} else if device.Type == models.Light {
				models.UpdateDeviceCols(&models.Device{
					DeviceID:   device.DeviceID,
					Brightness: ctx.ParamsInt64("value")}, "brightness")
			}
			break
		}
	}
}

func FaveHandler(ctx *macaron.Context) {
	for _, device := range models.GetDevices() {
		if device.DeviceID == ctx.ParamsInt64("id") {
			models.UpdateDeviceCols(&models.Device{
				DeviceID: device.DeviceID,
				IsFave:   !device.IsFave}, "is_fave")
			break
		}
	}
}

func RemoveHandler(ctx *macaron.Context) {
	models.DeleteDevice(ctx.ParamsInt64("id"))
}

func AddUserHandler(ctx *macaron.Context, form forms.RegisterForm, f *session.Flash) {
	ok := tokens.CheckAndConsumeKey(form.InviteCode)
	if !ok {
		f.Error("Invalid invite code. Please ask for an invite code from the home owner.")
		ctx.Redirect("/register")
		return
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
	if err != nil {
		panic(err)
	}
	user := &models.User{
		Username:  form.Email,
		Password:  string(pass),
		FirstName: form.FirstName,
		LastName:  form.LastName,
	}
	models.AddUser(user)
	ctx.Redirect("/dashboard")
}
