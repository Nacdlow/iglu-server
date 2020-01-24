package simulator

import (
	forms "gitlab.com/group-nacdlow/nacdlow-server/models/forms/sim"
	"gitlab.com/group-nacdlow/nacdlow-server/models/simulation"
	macaron "gopkg.in/macaron.v1"
	"time"
)

var (
	layoutCurrentTime = "January 2, 2006 3:04:05 PM"
	lastMCPing        time.Time
	mcVersion         string
)

// HomepageHandler handles rendering the simulator's homepage.
func HomepageHandler(ctx *macaron.Context) {
	ctx.Data["Env"] = simulation.Env
	ctx.Data["TickSleep"] = simulation.TickSleep.Milliseconds()
	ctx.Data["CurrentTime"] = time.Unix(simulation.Env.CurrentTime, 0).Format(layoutCurrentTime)
	if lastMCPing != (time.Time{}) {
		lastPing := time.Now().Sub(lastMCPing).Milliseconds()
		ctx.Data["LastMCPing"] = lastPing
		ctx.Data["MCConnected"] = (lastPing < 2000)
		ctx.Data["MCVersion"] = mcVersion
	}
	ctx.HTML(200, "sim/index")
}

// DataHandler handles posting the simulation data as JSON.
func DataHandler(ctx *macaron.Context) {
	if ctx.Query("from") == "minecraft" {
		lastMCPing = time.Now()
		mcVersion = ctx.Query("server")
	}
	ctx.JSON(200, simulation.Env)
}

func PostOverrideWeatherHandler(ctx *macaron.Context, form forms.OverrideWeatherForm) {
	simulation.Env.Weather.OutdoorTemp = form.OutdoorTemp
	simulation.Env.Weather.Humidity = form.Humidity
	simulation.Env.Weather.CloudCover = form.CloudCover
	ctx.Redirect("/sim")
}

func PostChangeTimeSleepHandler(ctx *macaron.Context, form forms.ChangeTimeSleepForm) {
	simulation.TickSleep = time.Duration(form.TickSleep) * time.Millisecond
	ctx.Redirect("/sim")
}
