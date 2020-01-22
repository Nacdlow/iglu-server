package simulation

import (
	"github.com/adlio/darksky"
)

// WeatherType represents the weather type.
type WeatherType int64

const (
	Clear = iota
	Cloudy
	Snow
	Rainy
)

// SimulationEnvironment represents an entire simulated environment state,
// which includes the home, weather, time and location states.
type SimulationEnvironment struct {
	Home
	Weather     WeatherStatus
	CurrentTime int64  // In Unix time.
	Location    string // Actual physical location.
}

// WeatherStatus represents a weather status state in the simulated
// environment.
type WeatherStatus struct {
	Type        WeatherType
	OutdoorTemp float64 // In Celcius.
	Humidity    float32 // In decimal, 0.5 = 50%.
	CloudCover  float32
	Alerts      []darksky.Alert
}

// Home represents a simulated home state.
type Home struct {
	MainDoorOpened bool // Whether the main door is opened or not.
	Rooms          []Room
	SolarMaxPower  int64 // Maximum solar panel generation capacity, in kWh.
}

// Room represents a simulated room state.
type Room struct {
	DBRoomID       int64 // The database room ID.
	Windows        []Window
	ActualRoomTemp float64
}

// Window represents a simulated window state.
type Window struct {
	RoomID int64
	IsOpen bool
}
