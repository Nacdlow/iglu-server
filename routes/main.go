package routes

import (
	"fmt"
	"html/template"
	"math"
	"strings"
	"time"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/models/forms"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/simulation"

	"github.com/go-macaron/session"
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
	ctx.Data["Stats"] = models.GetLatestStats()
	ctx.HTML(200, "dashboard")
}

// AlertsHandler handles the alerts page
func AlertsHandler(ctx *macaron.Context) {
	ctx.Data["CrossBack"] = 1
	ctx.Data["IsAlerts"] = 1
	ctx.HTML(200, "alerts")
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
