package routes

import (
	"net/http"
	"time"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/simulation"
	macaron "gopkg.in/macaron.v1"
)

// BatteryStatHandler handles rendering battery card.
func BatteryStatHandler(ctx *macaron.Context) {
	perc := int64((simulation.Env.BatteryStore / settings.Config.GetFloat64("Simulation.BatteryCapacityKWH")) * 100)
	if perc < 0 {
		perc = 0
	}
	if perc < 15 {
		ctx.Data["BatState"] = 0 // empty
	} else if perc < 50 {
		ctx.Data["BatState"] = 1 // quarter
	} else if perc < 75 {
		ctx.Data["BatState"] = 2 // half
	} else if perc < 95 {
		ctx.Data["BatState"] = 3 // three-quarter
	} else {
		ctx.Data["BatState"] = 4 // full
	}

	ctx.Data["BatteryPerc"] = perc
	ctx.HTML(http.StatusOK, "battery")
}

// ToggleHandler handles toggling devices.
func ToggleHandler(ctx *macaron.Context) {
	for _, device := range models.GetDevices() {
		if device.DeviceID == ctx.ParamsInt64("id") {
			if device.Type == models.Light || device.Type == models.Speaker ||
				device.Type == models.TempControl || device.Type == models.TV {
				err := models.UpdateDeviceCols(&models.Device{
					DeviceID:    device.DeviceID,
					Status:      !device.Status,
					ToggledUnix: time.Now().Unix()}, "status", "toggled_unix")
				if err != nil {
					panic(err)
				}
			}
			break
		}
	}
}

// SliderHandler handles modifying device values based on sliders.
func SliderHandler(ctx *macaron.Context) {
	for _, device := range models.GetDevices() {
		if device.DeviceID == ctx.ParamsInt64("id") {
			var err error
			if device.Type == models.Speaker {
				err = models.UpdateDeviceCols(&models.Device{
					DeviceID: device.DeviceID,
					Volume:   ctx.ParamsInt64("value")}, "volume")
			} else if device.Type == models.TempControl {
				err = models.UpdateDeviceCols(&models.Device{
					DeviceID: device.DeviceID,
					Temp:     ctx.ParamsFloat64("value")}, "temp")
			} else if device.Type == models.Light {
				err = models.UpdateDeviceCols(&models.Device{
					DeviceID:   device.DeviceID,
					Brightness: ctx.ParamsInt64("value")}, "brightness")
			}
			if err != nil {
				panic(err)
			}
			break
		}
	}
}

// FaveHandler handles favouriting or unfavouriting devices.
func FaveHandler(ctx *macaron.Context) {
	for _, device := range models.GetDevices() {
		if device.DeviceID == ctx.ParamsInt64("id") {
			err := models.UpdateDeviceCols(&models.Device{
				DeviceID: device.DeviceID,
				IsFave:   !device.IsFave}, "is_fave")
			if err != nil {
				panic(err)
			}
			break
		}
	}
}

// RemoveHandler handles removing devices.
func RemoveHandler(ctx *macaron.Context) {
	err := models.DeleteDevice(ctx.ParamsInt64("id"))
	if err != nil {
		panic(err)
	}
}

// RemoveRoomHandler handles removing rooms.
func RemoveRoomHandler(ctx *macaron.Context) {
	err := models.DeleteRoom(ctx.ParamsInt64("id"))
	if err != nil {
		panic(err)
	}
}
