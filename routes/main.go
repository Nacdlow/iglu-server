package routes

import (
	macaron "gopkg.in/macaron.v1"
)

// HomepageHandler handles rendering the homepage.
func HomepageHandler(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}

func SpecificRoomsHandler(ctx *macaron.Context) {
	ctx.HTML(200, "specificRooms")
}
//RoomsHandler handles rendering the rooms page
func RoomsHandler(ctx *macaron.Context) {
	ctx.HTML(200, "rooms")
}

func DevicesHandler(ctx *macaron.Context) {
	ctx.HTML(200, "devices")
}
