package routes

import (
	"fmt"
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/models/forms"
	"gitlab.com/group-nacdlow/nacdlow-server/models/simulation"
	macaron "gopkg.in/macaron.v1"
	"html/template"
	"math"
	"strings"
)

// NotFoundHandler handles 404 errors
func NotFoundHandler(ctx *macaron.Context) {
	ctx.HTML(404, "notfound")
}

// DashboardHandler handles rendering the dashboard.
func DashboardHandler(ctx *macaron.Context) {
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
	ctx.HTML(200, "dashboard")
}

// SpecificRoomsHandler handles the specific rooms
func SpecificRoomsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = fmt.Sprintf("%s", ctx.Params("roomType"))
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
	}
	ctx.Data["IsRooms"] = 1
	ctx.HTML(200, "specificRooms")
}

func AddDeviceRoomPostHandler(ctx *macaron.Context, form forms.AddDeviceForm) {
	device := &models.Device{
		RoomID:      form.RoomID,
		Type:        models.DeviceType(form.DeviceType),
		Description: form.Description,
	}
	if form.DeviceType == models.Light && form.IsMainLight {
		device.IsMainLight = true
	}
	models.AddDevice(device)
	ctx.Redirect(fmt.Sprintf("/room/%d", form.RoomID))
}

// RoomsHandler handles rendering the rooms page
func RoomsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Rooms"
	ctx.Data["IsRooms"] = 1
	ctx.Data["Rooms"] = models.GetRooms()
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

//LightsHandler handles the lights page
func LightsHandler(ctx *macaron.Context) {
	ctx.Data["Lights"] = models.GetDevices()
	ctx.HTML(200, "lights")
}

//HeatingHandler handles the temperature page
func HeatingHandler(ctx *macaron.Context) {
	ctx.Data["Heating"] = models.GetDevices()
	ctx.HTML(200, "temperature")
}

//SpeakerHandler handles the speakers page
func SpeakerHandler(ctx *macaron.Context) {
	ctx.Data["Speakers"] = models.GetDevices()
	ctx.HTML(200, "speakers")
}

func ToggleHandler(ctx *macaron.Context) {
	for _, device := range models.GetDevices() {
		if device.DeviceID == ctx.ParamsInt64("id") {
			if device.Type == models.Light || device.Type == models.Speaker ||
				device.Type == models.TempControl || device.Type == models.TV {
				models.UpdateDeviceCols(&models.Device{
					DeviceID: device.DeviceID,
					Status:   !device.Status}, "status")
			}
			break
		}
	}
}
