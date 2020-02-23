package routes

import (
	"net/http"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
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
	ctx.Data["CrossBack"] = 1
	ctx.Data["IsAdd"] = 1
	ctx.HTML(http.StatusOK, "add")
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

// SearchDeviceHandler handles the search for device
func SearchDeviceHandler(ctx *macaron.Context) {
	ctx.Data["BackLink"] = "/add"
	ctx.Data["CrossBack"] = 1
	ctx.Data["IsSearchDevice"] = 1
	ctx.HTML(http.StatusOK, "search_device")
}
