package simulator

import (
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	forms "gitlab.com/group-nacdlow/nacdlow-server/models/forms/sim"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/simulation"

	macaron "gopkg.in/macaron.v1"

	"net/http"
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
	if lastMCPing != (time.Time{}) {
		lastPing := time.Since(lastMCPing).Milliseconds()
		ctx.Data["LastMCPing"] = lastPing
		ctx.Data["MCConnected"] = (lastPing < 2000)
		ctx.Data["MCVersion"] = mcVersion
	}
	ctx.HTML(http.StatusOK, "sim/index")
}

// EnvStatusHandler handles loading the environment status for simulator page.
func EnvStatusHandler(ctx *macaron.Context) {
	ctx.Data["Env"] = simulation.Env
	ctx.Data["TickSleep"] = simulation.TickSleep.Milliseconds()
	ctx.Data["CurrentTime"] = time.Unix(simulation.Env.CurrentTime, 0).Format(layoutCurrentTime)
	ctx.HTML(http.StatusOK, "sim/env_status")
}

// DataHandler handles posting the simulation data as JSON.
func DataHandler(ctx *macaron.Context) {
	if ctx.Query("from") == "minecraft" {
		lastMCPing = time.Now()
		mcVersion = ctx.Query("server")
	}
	simulation.UpdateFromDB()
	ctx.JSON(200, simulation.Env)
}

// PostOverrideWeatherHandler handles post for overriding simulation's weather.
func PostOverrideWeatherHandler(ctx *macaron.Context, form forms.OverrideWeatherForm) {
	simulation.Env.Weather.OutdoorTemp = form.OutdoorTemp
	simulation.Env.Weather.Humidity = form.Humidity
	simulation.Env.Weather.CloudCover = form.CloudCover
	ctx.Redirect("/sim")
}

// PostChangeTimeSleepHandler handles post for changing simulation's tick speed.
func PostChangeTimeSleepHandler(ctx *macaron.Context, form forms.ChangeTimeSleepForm) {
	simulation.TickSleep = time.Duration(form.TickSleep) * time.Millisecond
	ctx.Redirect("/sim")
}

// ToggleHandler handles toggling a device in the simulation. This is made for
// the Minecraft simulation plugin.
func ToggleHandler(ctx *macaron.Context) {
	for _, room := range simulation.Env.Rooms {
		if room.MainLightDeviceID == ctx.ParamsInt64("id") {
			dev, err := models.GetDevice(room.MainLightDeviceID)
			if err == nil && dev.Type == models.Light {
				err := models.UpdateDeviceCols(&models.Device{
					DeviceID: dev.DeviceID,
					Status:   !dev.Status}, "status")
				if err != nil {
					panic(err)
				}
			}
			break
		}
	}
}

// SetMainDoorHandler sets the main door status.
func SetMainDoorHandler(ctx *macaron.Context) {
	simulation.Env.MainDoorOpened = (ctx.Params("status") == "opened")
}

// SetWindowStatusHandler sets the number of opened windows in a room.
func SetWindowStatusHandler(ctx *macaron.Context) {
	for _, room := range simulation.Env.Rooms {
		if room.DBRoomID == ctx.ParamsInt64("room") {
			room.OpenedWindows = ctx.ParamsInt("open_count")
			return
		}
	}
}
