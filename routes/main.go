package routes

import (
	macaron "gopkg.in/macaron.v1"
)

// HomepageHandler handles rendering the homepage.
func HomepageHandler(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}

// DashboardHandler handles rendering the dashboard.
func DashboardHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Dashboard"
	ctx.HTML(200, "dashboard")
}

func SpecificRoomsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "X's Room"
	ctx.HTML(200, "specificRooms")
}

//RoomsHandler handles rendering the rooms page
func RoomsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Rooms"
	ctx.HTML(200, "rooms")
}

func DevicesHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Devices"
	ctx.HTML(200, "devices")
}
