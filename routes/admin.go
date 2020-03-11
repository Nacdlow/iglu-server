package routes

import (
	"net/http"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/plugin"

	sdk "github.com/Nacdlow/plugin-sdk"
	macaron "gopkg.in/macaron.v1"
)

// AddDeviceHandler handles the add device page
func AddDeviceHandler(ctx *macaron.Context) {
	ctx.Data["BackLink"] = "/add"
	var err error
	ctx.Data["Rooms"], err = models.GetRooms()
	if err != nil {
		panic(err)
	}
	if ctx.Params("id") != "" {
		ctx.Data["RoomSelected"] = 1
		ctx.Data["RoomID"] = ctx.Params("id")
	}
	ctx.HTML(http.StatusOK, "add_device")
}

// AddHandler handles the add page
func AddHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Add..."
	ctx.Data["IsAdd"] = 1
	ctx.HTML(http.StatusOK, "add")
}

type PluginDeviceListing struct {
	PluginName string
	Devices    []AvailableDevice
}

type AvailableDevice struct {
	Device   sdk.AvailableDevice
	PluginID string
}

// SearchDeviceHandler handles the search for device
func SearchDeviceHandler(ctx *macaron.Context) {
	ctx.Data["BackLink"] = "/add"
	ctx.Data["CrossBack"] = 1
	ctx.Data["IsSearchDevice"] = 1
	ctx.HTML(http.StatusOK, "search_device")
}

// SearchDeviceListHandler handles the search for device
func SearchDeviceListHandler(ctx *macaron.Context) {
	var listings []PluginDeviceListing
	for _, plugin := range plugin.GetLoadedPlugins() {
		pluginDevices := plugin.Plugin.GetAvailableDevices()
		if len(pluginDevices) > 0 {
			var devices []AvailableDevice
			for _, d := range pluginDevices {
				if !models.HasPluginDevice(plugin.Manifest.Id, d.UniqueID) {
					devices = append(devices, AvailableDevice{Device: d, PluginID: plugin.Manifest.Id})
				}
			}

			listings = append(listings, PluginDeviceListing{
				PluginName: plugin.Manifest.Name,
				Devices:    devices,
			})
		}
	}
	ctx.Data["Listings"] = listings
	ctx.HTML(http.StatusOK, "search_device_list")
}

// AddRoomHandler handles the add room page
func AddRoomHandler(ctx *macaron.Context) {
	ctx.Data["BackLink"] = "/add"
	var err error
	ctx.Data["Rooms"], err = models.GetRooms()
	if err != nil {
		panic(err)
	}
	ctx.HTML(http.StatusOK, "add_room")
}

// AddScheduleHandler handles the adding schedule page
func AddScheduleHandler(ctx *macaron.Context) {
	ctx.Data["BackLink"] = "/add"

	ctx.HTML(http.StatusOK, "add_schedule")
}
