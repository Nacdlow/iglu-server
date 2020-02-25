package routes

import (
	"net/http"
	"time"

	"github.com/go-macaron/session"
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
	devices, err := models.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
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
	devices, err := models.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		if device.DeviceID == ctx.ParamsInt64("id") {
			var err error
			switch device.Type {
			case models.Speaker:
				err = models.UpdateDeviceCols(&models.Device{
					DeviceID: device.DeviceID,
					Volume:   ctx.ParamsInt64("value")}, "volume")
			case models.TempControl:
				err = models.UpdateDeviceCols(&models.Device{
					DeviceID: device.DeviceID,
					Temp:     ctx.ParamsFloat64("value")}, "temp")
			case models.Light:
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

// FontSliderHandler handles modifying font size values based on sliders.
func FontSliderHandler(ctx *macaron.Context, sess session.Store) {
	if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
		var fontSize string

		switch ctx.ParamsInt("size") {
		case 0:
			fontSize = "font-xsmall"
		case 1:
			fontSize = "font-small"
		case 2:
			fontSize = "font-medium"
		case 3:
			fontSize = "font-large"
		case 4:
			fontSize = "font-xlarge"
		}

		models.UpdateUserCols(&models.User{Username: user.Username, FontSize: fontSize}, "font_size")
	}
}

// FaveHandler handles favouriting or unfavouriting devices.
func FaveHandler(ctx *macaron.Context) {
	devices, err := models.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
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

func ChangeNameHandler(ctx *macaron.Context) {
	rooms, err := models.GetRooms()
	if err != nil {
		panic(err)
	}
	for _, room := range rooms {
		if room.RoomID == ctx.ParamsInt64("id") {
			err := models.UpdateRoomCols(&models.Room{
				RoomID:   room.RoomID,
				RoomName: ctx.Params("newName"),
			})
			if err != nil {
				panic(err)
			}
			break
		}
	}
}
