package routes

import (
	"time"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/simulation"
	macaron "gopkg.in/macaron.v1"
)

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
	ctx.HTML(200, "battery")
}

func ToggleHandler(ctx *macaron.Context) {
	for _, device := range models.GetDevices() {
		if device.DeviceID == ctx.ParamsInt64("id") {
			if device.Type == models.Light || device.Type == models.Speaker ||
				device.Type == models.TempControl || device.Type == models.TV {
				models.UpdateDeviceCols(&models.Device{
					DeviceID:    device.DeviceID,
					Status:      !device.Status,
					ToggledUnix: time.Now().Unix()}, "status", "toggled_unix")
			}
			break
		}
	}
}

func SliderHandler(ctx *macaron.Context) {
	for _, device := range models.GetDevices() {
		if device.DeviceID == ctx.ParamsInt64("id") {
			if device.Type == models.Speaker {
				models.UpdateDeviceCols(&models.Device{
					DeviceID: device.DeviceID,
					Volume:   ctx.ParamsInt64("value")}, "volume")
			} else if device.Type == models.TempControl {
				models.UpdateDeviceCols(&models.Device{
					DeviceID: device.DeviceID,
					Temp:     ctx.ParamsFloat64("value")}, "temp")
			} else if device.Type == models.Light {
				models.UpdateDeviceCols(&models.Device{
					DeviceID:   device.DeviceID,
					Brightness: ctx.ParamsInt64("value")}, "brightness")
			}
			break
		}
	}
}

func FaveHandler(ctx *macaron.Context) {
	for _, device := range models.GetDevices() {
		if device.DeviceID == ctx.ParamsInt64("id") {
			models.UpdateDeviceCols(&models.Device{
				DeviceID: device.DeviceID,
				IsFave:   !device.IsFave}, "is_fave")
			break
		}
	}
}

func RemoveHandler(ctx *macaron.Context) {
	models.DeleteDevice(ctx.ParamsInt64("id"))
}

func RemoveRoomHandler(ctx *macaron.Context) {
	models.DeleteRoom(ctx.ParamsInt64("id"))
}
