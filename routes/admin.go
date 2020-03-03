package routes

import (
	"net/http"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/plugin"

	"github.com/Nacdlow/plugin-sdk"
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
	Devices    []sdk.AvailableDevice
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
		devices := plugin.Plugin.GetAvailableDevices()
		if len(devices) > 0 {
			listings = append(listings, PluginDeviceListing{
				PluginName: plugin.Plugin.GetManifest().Name,
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
