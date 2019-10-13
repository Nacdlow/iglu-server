package routes

import (
	macaron "gopkg.in/macaron.v1"
)

// HomepageHandler handles rendering the homepage.
func HomepageHandler(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}
