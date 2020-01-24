package simulation

import (
	"github.com/adlio/darksky"
	log "github.com/sirupsen/logrus"
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/weather"
	"time"
)

var (
	// Env is the current simulation Environment
	Env Environment
	// TickSleep is the length of a second in the simulation world, in
	// milliseconds.
	TickSleep time.Duration = 1000
)

// LoadFromDB loads the rooms from the database into
func LoadFromDB() {
	for _, room := range models.GetRooms() {
		r := Room{
			DBRoomID:    room.RoomID,
			LightStatus: false, // Assume lights are off
		}

		// Add Windows
		for i := int64(0); i < room.WindowCount; i++ {
			r.Windows = append(r.Windows, Window{false}) // Assume all windows are closed
		}

		// Get room temp and light status from the devices of that room
		tempSet, lightSet := false, false
		for _, device := range models.GetDevices() {
			if device.RoomID == room.RoomID {
				switch device.Type {
				case models.TempControl:
					r.ActualRoomTemp = device.Temp
					r.TempControlDeviceID = device.DeviceID
					tempSet = true
				case models.Light:
					if device.IsMainLight {
						r.LightStatus = device.Status
						r.MainLightDeviceID = device.DeviceID
						lightSet = true
					}
				}
			}
		}

		if !tempSet {
			if Env.ForecastData != nil {
				log.Errorf("Room %d (%s) does not have a temperature control device! Using outside temp...",
					room.RoomID, room.RoomName)
				r.ActualRoomTemp = Env.ForecastData.Currently.Temperature
			} else {
				log.Errorf("Room %d (%s) does not have a temperature control device! Cannot use outside temp either! Setting to 20C.",
					room.RoomID, room.RoomName)
				r.ActualRoomTemp = 20
			}
		}
		if !lightSet {
			log.Errorf("Room %d (%s) does not have a main light source (device)!",
				room.RoomID, room.RoomName)
		}

		Env.Rooms = append(Env.Rooms, r)
	}
}

// Start will start the simulation environment and make it Tick at TickSpeed.
// This function will not return (halt).
func Start() {
	// Load current time for the simulation
	Env.CurrentTime = time.Now().Unix()
	// Load the forecast data
	f, err := weather.GetWeather()
	if err != nil {
		log.Error("Failed to get forecast!", err)
		log.Error("Please make sure that the Darksky API key is correct.")
		log.Error("You may use the group API key at: https://wiki.nacdlow.com/Accounts.html")
		log.Error("Simulation will may not work properly!")
	}
	Env.ForecastData = f
	LoadFromDB()
	for {
		Tick()
		time.Sleep(TickSleep * time.Millisecond)
	}
}

// Tick will tick the environment one second.
func Tick() {
	Env.CurrentTime++
	// TODO update room temperatures
	for _, room := range Env.Home.Rooms {
		for _, window := range room.Windows {
			if window.IsOpen {
				// TODO get the difference between outdoor temp and room
			}
		}
	}
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
	Home         `json:"home"`
	Weather      WeatherStatus     `json:"weather"`
	CurrentTime  int64             `json:"current_time"` // In Unix time.
	ForecastData *darksky.Forecast `json:"-"`
}

// WeatherStatus represents a weather status state in the simulated
// environment.
type WeatherStatus struct {
	Type        WeatherType `json:"type"`
	OutdoorTemp float64     `json:"outdoor_temp"` // In Celcius.
	Humidity    float32     `json:"humidity"`     // In decimal, 0.5 = 50%.
	CloudCover  float32     `json:"cloud_cover"`
}

// Home represents a simulated home state.
type Home struct {
	MainDoorOpened bool   `json:"main_door_opened"` // Whether the main door is opened or not.
	Rooms          []Room `json:"rooms"`
	SolarMaxPower  int64  `json:"solar_max_power"` // Maximum solar panel generation capacity, in kWh.
}

// Room represents a simulated room state.
type Room struct {
	DBRoomID            int64    `json:"db_room_id"` // The database room ID.
	Windows             []Window `json:"windows"`
	ActualRoomTemp      float64  `json:"actual_room_temp"`
	TempControlDeviceID int64    `json:"temp_control_device_id"`
	LightStatus         bool     `json:"light_status"`
	MainLightDeviceID   int64    `json:"main_light_device_id"`
}

// Window represents a simulated window state.
type Window struct {
	IsOpen bool `json:"is_open"`
}
