package routes

import (
	"fmt"
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	macaron "gopkg.in/macaron.v1"
)

// NotFoundHandler handles 404 errors
func NotFoundHandler(ctx *macaron.Context) {
	ctx.HTML(404, "notfound")
}

// HomepageHandler handles rendering the homepage.
func HomepageHandler(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}

// DashboardHandler handles rendering the dashboard.
func DashboardHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Dashboard"
	ctx.Data["IsDashboard"] = 1
	ctx.HTML(200, "dashboard")
}

// SpecificRoomsHandler handles the specific rooms
func SpecificRoomsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = fmt.Sprintf("%s", ctx.Params("roomType"))
	ctx.Data["IsRooms"] = 1
	ctx.HTML(200, "specificRooms")
}

// RoomsHandler handles rendering the rooms page
func RoomsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Rooms"
	ctx.Data["IsRooms"] = 1
	ctx.Data["Rooms"] = models.GetRooms()
	ctx.HTML(200, "rooms")
}

// AddRoomHandler handles rendering the add room page page.
func AddRoomHandler(ctx *macaron.Context) {
	ctx.HTML(200, "addroom")
}

// DevicesHandler handles the devices page
func DevicesHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Devices"
	ctx.HTML(200, "devices")
}

// RegisterHandler handles the registration page.
func RegisterHandler(ctx *macaron.Context) {
	ctx.HTML(200, "register")
}

//LightsHandler handles the lights page
func LightsHandler(ctx *macaron.Context) {
	ctx.HTML(200, "lights")
}
