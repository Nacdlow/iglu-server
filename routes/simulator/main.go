package simulator

import (
	macaron "gopkg.in/macaron.v1"
)

// HomepageHandler handles rendering the simulator's homepage.
func HomepageHandler(ctx *macaron.Context) {
	ctx.HTML(200, "sim/index")
}
