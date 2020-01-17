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
	ctx.HTML(200, "dashboard")
}

func DevicesHandler(ctx *macaron.Context) {
	ctx.HTML(200, "devices")
}
