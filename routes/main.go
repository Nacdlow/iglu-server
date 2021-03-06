package routes

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Nacdlow/iglu-server/models"
	"github.com/Nacdlow/iglu-server/models/forms"
	"github.com/Nacdlow/iglu-server/modules/plugin"
	"github.com/Nacdlow/iglu-server/modules/simulation"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

// NotFoundHandler handles 404 errors
func NotFoundHandler(ctx *macaron.Context) {
	ctx.HTML(http.StatusNotFound, "notfound")
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
	var err error
	ctx.Data["Devices"], err = models.GetDevices()
	if err != nil {
		panic(err)
	}
	ctx.Data["Stats"], err = models.GetLatestStats()
	if err != nil {
		panic(err)
	}
	ctx.HTML(http.StatusOK, "dashboard")
}

// AlertsHandler handles the alerts page
func AlertsHandler(ctx *macaron.Context, sess session.Store) {
	ctx.Data["CrossBack"] = 1
	ctx.Data["IsAlerts"] = 1
	alerts, err := models.GetUserAlerts(sess.Get("username").(string))
	if err != nil {
		panic(err)
	}
	ctx.Data["None"] = (len(alerts) == 0)
	ctx.Data["Alerts"] = alerts
	ctx.HTML(http.StatusOK, "alerts")
}

// InternalAccounts handels the internal accounts page
func InternalAccountsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Accounts"
	var err error
	ctx.Data["Accounts"], err = models.GetUsers()
	if err != nil {
		panic(err)
	}
	ctx.Data["ArrowBack"] = 1
	ctx.HTML(http.StatusNotFound, "internal_accounts")
}

// SpecificRoomsHandler handles the specific rooms
func SpecificRoomsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = ctx.Params("roomType")
	rooms, err := models.GetRooms()
	if err != nil {
		panic(err)
	}
	if ctx.Params("name") == "bedrooms" {
		ctx.Data["Bedrooms"] = 1
		ctx.Data["Rooms"] = rooms
	} else if ctx.Params("name") == "bathrooms" {
		ctx.Data["Bathrooms"] = 1
		ctx.Data["Rooms"] = rooms
	} else {
		room, err := models.GetRoom(ctx.ParamsInt64("name"))
		if err != nil {
			ctx.Redirect("/rooms")
			return
		}
		ctx.Data["NavTitle"] = room.RoomName
		ctx.Data["Room"] = room
		ctx.Data["Rooms"], err = models.GetRooms()
		ctx.Data["Devices"], err = models.GetDevices()
		if err != nil {
			panic(err)
		}
	}

	ctx.Data["ArrowBack"] = 1
	ctx.Data["IsRooms"] = 1
	ctx.Data["IsSpecificRoom"] = 1
	ctx.HTML(http.StatusOK, "specificRooms")
}

// OverviewHandler handles the overview page
func OverviewHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Overview"
	ctx.Data["IsOverview"] = 1
	var err error
	ctx.Data["Devices"], err = models.GetDevices()
	if err != nil {
		panic(err)
	}
	ctx.HTML(http.StatusOK, "overview")
}

// AddDeviceRoomPostHandler handles post for adding a device to a room.
func AddDeviceRoomPostHandler(ctx *macaron.Context, form forms.AddDeviceForm,
	errs binding.Errors, f *session.Flash) {
	if len(errs) > 0 {
		f.Error("Missing required fields!")
		ctx.Redirect(fmt.Sprintf("/rooms"))
		return
	}
	device := &models.Device{
		RoomID:      form.RoomID,
		Type:        models.DeviceType(form.DeviceType),
		Description: form.Description,
		ToggledUnix: time.Now().Unix(),
	}
	devices, err := models.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, d := range devices {
		if d.RoomID == device.RoomID && d.Type == device.Type {
			switch device.Type {
			case models.Light: // We can't have two main lights
				if form.IsMainLight && d.IsMainLight {
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
	err = models.AddDevice(device)
	if err != nil {
		panic(err)
	}
	ctx.Redirect(fmt.Sprintf("/room/%d", form.RoomID))
}

var TYPE_NAME = []string{"Light", "Temp Control", "Other", "Speaker"}

func ConnectDeviceHandler(ctx *macaron.Context, f *session.Flash) {
	pl, err := plugin.GetPlugin(ctx.Params("plugin"))
	if err != nil {
		f.Error("Cannot connect to plugin")
		ctx.Redirect("/rooms")
		return
	}
	devices := pl.Plugin.GetAvailableDevices()
	log.Println(devices)
	for _, dev := range devices {
		if dev.UniqueID == ctx.Params("id") {
			ctx.Data["Device"] = dev
			var err error
			ctx.Data["Type"] = TYPE_NAME[dev.Type]
			if err != nil {
				panic(err)
			}
			rooms, err := models.GetRooms()
			if err != nil {
				panic(err)
			}
			ctx.Data["Rooms"] = rooms
			ctx.HTML(http.StatusOK, "connect_device")
			return
		}
	}
	f.Error("Cannot find device (doesn't exist anymore?)")
	ctx.Redirect("/rooms")
}

func IdentifyDeviceHandler(ctx *macaron.Context, f *session.Flash) {
	pl, err := plugin.GetPlugin(ctx.Params("plugin"))
	if err != nil {
		return
	}

	devices := pl.Plugin.GetAvailableDevices()
	for _, dev := range devices {
		if dev.UniqueID == ctx.Params("id") {
			go func() {
				for i := 0; i < 3; i++ {
					pl.Plugin.OnDeviceToggle(ctx.Params("id"), true)
					time.Sleep(800 * time.Millisecond)
					pl.Plugin.OnDeviceToggle(ctx.Params("id"), false)
					time.Sleep(800 * time.Millisecond)
				}
			}()
			return
		}
	}
}

func ConnectDevicePostHandler(ctx *macaron.Context, form forms.AddDeviceForm,
	errs binding.Errors, f *session.Flash) {
	if len(errs) > 0 {
		f.Error("Missing required fields!")
		ctx.Redirect("/rooms")
		return
	}
	pl, err := plugin.GetPlugin(ctx.Params("plugin"))
	if err != nil {
		f.Error("Cannot connect to plugin")
		ctx.Redirect("/rooms")
		return
	}
	devices := pl.Plugin.GetAvailableDevices()
	log.Println(devices)
	for _, dev := range devices {
		if dev.UniqueID == ctx.Params("id") {
			device := &models.Device{
				RoomID:         form.RoomID,
				Type:           models.DeviceType(form.DeviceType),
				Description:    form.Description,
				ToggledUnix:    time.Now().Unix(),
				IsRegistered:   true,
				PluginID:       ctx.Params("plugin"),
				PluginUniqueID: ctx.Params("id"),
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
			err = models.AddDevice(device)
			if err != nil {
				panic(err)
			}
			ctx.Redirect(fmt.Sprintf("/room/%d", form.RoomID))
			return
		}
	}
	f.Error("Cannot find device (doesn't exist anymore?)")
	ctx.Redirect("/rooms")
}

// RoomsHandler handles rendering the rooms page
func RoomsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Rooms"
	ctx.Data["IsRooms"] = 1
	rooms, err := models.GetRooms()
	if err != nil {
		panic(err)
	}
	for i := range rooms {
		rooms[i].LoadMainDevices()
	}
	ctx.Data["Rooms"] = rooms
	ctx.HTML(http.StatusOK, "rooms")
}

// PostRoomHandler handles post request for room page, to add a room.
func PostRoomHandler(ctx *macaron.Context, form forms.AddRoomForm,
	errs binding.Errors, f *session.Flash) {
	if len(errs) > 0 {
		f.Error("Missing required fields!")
		ctx.Redirect("/rooms")
		return
	}
	room := &models.Room{
		RoomName:     form.RoomName,
		RoomType:     models.RType(form.RoomType),
		IsRestricted: form.IsRestricted,
	}
	if form.PartOfRoom >= 0 {
		room.IsSubRoom = true
		room.PartOfRoom = form.PartOfRoom
	}
	err := models.AddRoom(room)
	if err != nil {
		panic(err)
	}
	ctx.Redirect("/rooms")
}

// DevicesHandler handles the devices page
func DevicesHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Devices"
	ctx.HTML(http.StatusOK, "devices")
}
