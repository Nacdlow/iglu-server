package simulation

import (
	"github.com/adlio/darksky"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
	"log"
)

var Env Environment

// Tick will tick the environment one second.
func Tick() {
	Env.CurrentTime++
	// TODO update room temperatures
}

// GetWeather loads the forecast into the simulation.
func GetWeather() {
	client := darksky.NewClient(settings.Config.GetString("DarkskyAPIKey"))
	f, err := client.GetForecast(Env.Location.Latitude, Env.Location.Longitude,
		darksky.Defaults)
	if err != nil {
		log.Fatal("Failed to get forecast! ", err)
		return
	}
	Env.ForecastData = f
}

// WeatherType represents the weather type.
type WeatherType int64

// WeatherType enums.
const (
	Clear = iota
	Cloudy
	Snow
	Rainy
)

// Environment represents an entire simulated environment state, which includes
// the home, weather, time and location states.
type Environment struct {
	Home        `json:"home"`
	Weather     WeatherStatus `json:"weather"`
	CurrentTime int64         `json:"current_time"` // In Unix time.
	Location    struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"location"`
	ForecastData *darksky.Forecast `json:"-"`
}

// WeatherStatus represents a weather status state in the simulated
// environment.
type WeatherStatus struct {
	Type        WeatherType     `json:"type"`
	OutdoorTemp float64         `json:"outdoor_temp"` // In Celcius.
	Humidity    float32         `json:"humidity"`     // In decimal, 0.5 = 50%.
	CloudCover  float32         `json:"cloud_cover"`
	Alerts      []darksky.Alert `json:"alerts"`
}

// Home represents a simulated home state.
type Home struct {
	MainDoorOpened bool   `json:"main_door_opened"` // Whether the main door is opened or not.
	Rooms          []Room `json:"rooms"`
	SolarMaxPower  int64  `json:"solar_max_power"` // Maximum solar panel generation capacity, in kWh.
}

// Room represents a simulated room state.
type Room struct {
	DBRoomID       int64    `json:"db_room_id"` // The database room ID.
	Windows        []Window `json:"windows"`
	ActualRoomTemp float64  `json:"actual_room_temp"`
}

// Window represents a simulated window state.
type Window struct {
	RoomID int64 `json:"room_id"`
	IsOpen bool  `json:"is_open"`
}
