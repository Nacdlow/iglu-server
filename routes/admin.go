package routes

import (
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	macaron "gopkg.in/macaron.v1"
)

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
