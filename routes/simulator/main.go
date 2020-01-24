package simulator

import (
	"gitlab.com/group-nacdlow/nacdlow-server/models/simulation"
	macaron "gopkg.in/macaron.v1"
	"time"
)

var (
	layoutCurrentTime = "January 2, 2006 3:04:05 PM"
)

// HomepageHandler handles rendering the simulator's homepage.
func HomepageHandler(ctx *macaron.Context) {
	ctx.Data["Env"] = simulation.Env
	ctx.Data["CurrentTime"] = time.Unix(simulation.Env.CurrentTime, 0).Format(layoutCurrentTime)
	ctx.HTML(200, "sim/index")
}

// DataHandler handles posting the simulation data as JSON.
func DataHandler(ctx *macaron.Context) {
	ctx.JSON(200, simulation.Env)
}
