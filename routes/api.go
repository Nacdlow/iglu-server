package routes

import (
	"net/http"
	"time"

	"github.com/go-macaron/session"
	"github.com/Nacdlow/iglu-server/models"
	"github.com/Nacdlow/iglu-server/modules/plugin"
	"github.com/Nacdlow/iglu-server/modules/settings"
	"github.com/Nacdlow/iglu-server/modules/simulation"
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
func ToggleHandler(ctx *macaron.Context, sess session.Store) {
	devices, err := models.GetDevices()
	if err != nil {
		panic(err)
	}

	for _, device := range devices {
		if device.DeviceID == ctx.ParamsInt64("id") {
			if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
				room, err := models.GetRoom(device.RoomID)
				if err != nil {
					panic(err)
				}
				if user.Role != models.AdminRole && room.IsRestricted {
					return
				}
			}
			if device.Type == models.Light || device.Type == models.Speaker ||
				device.Type == models.TempControl || device.Type == models.Other {
				err := models.UpdateDeviceCols(&models.Device{
					DeviceID:    device.DeviceID,
					Status:      !device.Status,
					ToggledUnix: time.Now().Unix()}, "status", "toggled_unix")
				if err != nil {
					panic(err)
				}
				if device.IsRegistered {
					pl, err := plugin.GetPlugin(device.PluginID)
					if err != nil {
						panic(err)
					}
					err = pl.Plugin.OnDeviceToggle(device.PluginUniqueID, !device.Status)
					if err != nil {
						panic(err)
					}
				}
			}
			break
		}
	}
}

// SliderHandler handles modifying device values based on sliders.
func SliderHandler(ctx *macaron.Context, sess session.Store) {
	devices, err := models.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		if device.DeviceID == ctx.ParamsInt64("id") {
			if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
				room, err := models.GetRoom(device.RoomID)
				if err != nil {
					panic(err)
				}
				if user.Role != models.AdminRole && room.IsRestricted {
					return
				}
			}
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
			fontSize = "xsmall"
		case 1:
			fontSize = "small"
		case 2:
			fontSize = "medium"
		case 3:
			fontSize = "large"
		case 4:
			fontSize = "xlarge"
		}

		models.UpdateUserCols(&models.User{Username: user.Username, FontSize: fontSize}, "font_size")
	}
}

// FontDropDownHandler handles changing the font typeface
func FontDropDownHandler(ctx *macaron.Context, sess session.Store) {
	if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
		var fontFace string

		switch ctx.ParamsInt("font") {
		case 0:
			fontFace = "Roboto"
		case 1:
			fontFace = "sans-serif"
		case 2:
			fontFace = "Roboto Light"
		case 3:
			fontFace = "Roboto Bold"
		case 4:
			fontFace = "OpenDyslexic"
		case 5:
			fontFace = "OpenDyslexic Bold"
		case 6:
			fontFace = "OpenDyslexic3"
		case 7:
			fontFace = "OpenDyslexic3 Bold"
		}

		models.UpdateUserCols(&models.User{Username: user.Username, FontFace: fontFace}, "font_face")

		ctx.Redirect("/settings/appearance")
	}
}

// FaveHandler handles favouriting or unfavouriting devices.
func FaveHandler(ctx *macaron.Context, sess session.Store) {
	devices, err := models.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		if device.DeviceID == ctx.ParamsInt64("id") {
			if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
				room, err := models.GetRoom(device.RoomID)
				if err != nil {
					panic(err)
				}
				if user.Role != models.AdminRole && room.IsRestricted {
					return
				}
			}
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

func RestrictHandler(ctx *macaron.Context) {
	rooms, err := models.GetRooms()
	if err != nil {
		panic(err)
	}
	for _, room := range rooms {
		if room.RoomID == ctx.ParamsInt64("id") {
			err := models.UpdateRoomCols(&models.Room{
				RoomID:       room.RoomID,
				IsRestricted: !room.IsRestricted}, "is_restricted")
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
func RemoveRoomHandler(ctx *macaron.Context, sess session.Store) {
	if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
		room, err := models.GetRoom(ctx.ParamsInt64("id"))
		if err != nil {
			panic(err)
		}
		if user.Role != models.AdminRole && room.IsRestricted {
			return
		}
	}
	err := models.DeleteRoom(ctx.ParamsInt64("id"))
	if err != nil {
		panic(err)
	}
}

// ChangeNameHandler add comments @Ruaridh!
func ChangeNameHandler(ctx *macaron.Context, sess session.Store) {
	rooms, err := models.GetRooms()
	if err != nil {
		panic(err)
	}
	for _, room := range rooms {
		if room.RoomID == ctx.ParamsInt64("id") {
			if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
				if user.Role != models.AdminRole && room.IsRestricted {
					return
				}
			}
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

//ChangeDeviceNameHandler handles changing of the device name, crazy huh?
func ChangeDeviceNameHandler(ctx *macaron.Context, sess session.Store) {
	devices, err := models.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		if device.DeviceID == ctx.ParamsInt64("id") {
			if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
				room, err := models.GetRoom(device.RoomID)
				if err != nil {
					panic(err)
				}
				if user.Role != models.AdminRole && room.IsRestricted {
					return
				}
			}
			err := models.UpdateDeviceCols(&models.Device{
				DeviceID:    device.DeviceID,
				Description: ctx.Params("newName"),
			})
			if err != nil {
				panic(err)
			}
			break
		}
	}
}

func MoveDeviceHandler(ctx *macaron.Context, sess session.Store) {
	devices, err := models.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		if device.DeviceID == ctx.ParamsInt64("did") {
			if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
				room, err := models.GetRoom(device.RoomID)
				if err != nil {
					panic(err)
				}
				if user.Role != models.AdminRole && room.IsRestricted {
					return
				}
			}
			err := models.UpdateDeviceCols(&models.Device{
				DeviceID: device.DeviceID,
				RoomID:   ctx.ParamsInt64("rid"),
			})
			if err != nil {
				panic(err)
			}
			break
		}
	}
}
