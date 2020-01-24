package simulator

import (
	"gitlab.com/group-nacdlow/nacdlow-server/models/simulation"
	macaron "gopkg.in/macaron.v1"
)

// HomepageHandler handles rendering the simulator's homepage.
func HomepageHandler(ctx *macaron.Context) {
	ctx.HTML(200, "sim/index")
}

// DataHandler handles posting the simulation data as JSON.
func DataHandler(ctx *macaron.Context) {
	ctx.JSON(200, simulation.Env)
}
