package simulation

import (
	"github.com/adlio/darksky"
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/weather"
	"log"
	"math"
	"math/rand"
	"time"
)

var (
	// Env is the current simulation Environment
	Env Environment
	// TickSleep is the length of a second in the simulation world, in
	// milliseconds.
	TickSleep time.Duration = 1000 * time.Millisecond
	rnd       *rand.Rand
)

func init() {
	// Create and seed a new randomiser
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func UpdateFromDB() {
	for i, room := range Env.Rooms {
		device, err := models.GetDevice(room.MainLightDeviceID)
		if err == nil && device.Type == models.Light {
			if device.Status && device.Brightness < 3 {
				Env.Rooms[i].LightStatus = false
			} else {
				Env.Rooms[i].LightStatus = device.Status
			}
		}
	}
}

// LoadFromDB loads the rooms from the database into
func LoadFromDB() {
	lacking := false
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
						r.MainLightDeviceID = device.DeviceID
						lightSet = true
						// brightness threshold (minecraft has no brightness for lights
						if device.Status && device.Brightness < 3 {
							r.LightStatus = false
						} else {
							r.LightStatus = device.Status
						}
					}
				}
			}
		}

		// If the temperature isn't set (due to no temp control), we can set it
		// to something else.
		if !tempSet {
			if Env.ForecastData != nil {
				log.Printf("Room %d (%s) does not have a temperature control device! Using outside temp...\n",
					room.RoomID, room.RoomName)
				r.ActualRoomTemp = Env.ForecastData.Currently.Temperature
			} else {
				log.Printf("Room %d (%s) does not have a temperature control device! Cannot use outside temp either! Setting to 20C.\n",
					room.RoomID, room.RoomName)
				r.ActualRoomTemp = 20
			}
		}

		// No room light source? Be angry!
		if !lightSet {
			log.Printf("Room %d (%s) does not have a main light source (device)!\n",
				room.RoomID, room.RoomName)
		}

		if !lightSet || !tempSet {
			lacking = true
		}

		Env.Rooms = append(Env.Rooms, r)
	}
	if lacking {
		log.Println("Some rooms may be lacking important devices which is required for the simulator to work properly!" +
			" Please fix before continuing.")
	}
}

// Start will start the simulation environment and make it Tick at TickSpeed.
// This function will not return (halt).
func Start() {
	Env.SolarMaxPower = settings.Config.GetInt("Simulation.SolarCapacityKW")
	// Load current time for the simulation
	Env.CurrentTime = time.Now().Unix()
	// Load the forecast data
	f, err := weather.GetWeather()
	if err != nil {
		log.Println("Failed to get forecast! ", err)
		log.Println("Please make sure that the Darksky API key is correct.")
		log.Println("You may use the group API key at: https://wiki.nacdlow.com/Accounts.html")
		log.Println("Simulation may not work properly!")
		Env.Weather.OutdoorTemp = 4
		Env.Weather.Humidity = 0.12
		Env.Weather.CloudCover = 0.45
	} else {
		Env.ForecastData = f
		current := f.Currently
		Env.Weather.OutdoorTemp = current.Temperature
		Env.Weather.Humidity = current.Humidity
		Env.Weather.CloudCover = current.CloudCover
	}
	// Load the simulation environment from DB
	LoadFromDB()
	// This is the main simulation loop
	for {
		Tick()
		time.Sleep(TickSleep)
	}
}

// getChange makes the current go closer to the influence temperature depending
// on the change. The higher the change value is, the less (slower) it goes
// towards to influence.
//
// This is used for making the room temperature move towards an "influence",
// whether that'd be the air conditioning or the outside temperature (if the
// window is opened).
func getChange(current, influence, change, threshold float64) float64 {
	diff := math.Abs(current - influence)
	if diff <= threshold { // threshold
		return influence
	} else {
		changed := diff / change
		if influence > current {
			return current + changed
		} else {
			return current - changed
		}
	}
}

// Tick will tick the simulated environment one second.
func Tick() {
	Env.CurrentTime++
	var runningTempCont, runningLights int
	// Update room temperatures
	UpdateFromDB()
	outTemp := Env.Weather.OutdoorTemp
	for i, room := range Env.Home.Rooms {
		// Simulate cold/hot air from outside coming in room through windows
		for _, window := range room.Windows {
			if window.IsOpen {
				Env.Home.Rooms[i].ActualRoomTemp = getChange(room.ActualRoomTemp, outTemp, 270, 0.75)
			}
		}

		if room.LightStatus {
			runningLights++
		}

		// Simulate the room heating/cooling from outside temperature "leak"
		Env.Home.Rooms[i].ActualRoomTemp = getChange(room.ActualRoomTemp, outTemp, 400, 0)

		// Simulate the temperature control heating/cooling the room
		tempCont, err := models.GetDevice(room.TempControlDeviceID)
		if err == nil && tempCont.Type == models.TempControl && tempCont.Status {
			runningTempCont++
			Env.Home.Rooms[i].ActualRoomTemp = getChange(room.ActualRoomTemp, tempCont.Temp, 240, 0.75)
		}
	}
	now := time.Unix(Env.CurrentTime, 0)

	// Calculate power consumption
	powerPerTempCont := (float64(15) / float64(len(Env.Home.Rooms))) // 15kW per entire house
	Env.PowerConRate = 11                                            // Baseline kW consumption
	conVary := math.Sin(float64(now.Second())/10) / 6
	Env.PowerConRate += conVary
	Env.PowerConRate += (float64(runningLights) * 0.20) // Light power consumption
	Env.PowerConRate += (powerPerTempCont * float64(runningTempCont))

	Env.NetPower = Env.PowerGenRate - Env.PowerConRate

	Env.BatteryStore += Env.NetPower / 60 / 60
	if Env.BatteryStore > settings.Config.GetFloat64("Simulation.BatteryCapacityKWH") {
		Env.BatteryStore = settings.Config.GetFloat64("Simulation.BatteryCapacityKWH")
	}

	// Update kWh solar generation value
	var change float64
	if Env.Weather.CloudCover > 0 {
		change = float64(Env.SolarMaxPower) * (1 - Env.Weather.CloudCover)
	} else {
		change = float64(Env.SolarMaxPower)
	}

	change += math.Sin(float64(now.Second())/10) * 5
	if change < 0 {
		Env.PowerGenRate = 0
	} else if change > float64(Env.SolarMaxPower) {
		Env.PowerGenRate = float64(Env.SolarMaxPower)
	} else {
		Env.PowerGenRate = change
	}
	if now.Hour() < 6 || now.Hour() > 18 { // Night
		Env.PowerGenRate = 0
	} else if now.Hour() < 9 || now.Hour() > 17 { // Early morning/late afternoon
		Env.PowerGenRate = Env.PowerGenRate / 4
	}

	Env.MinecraftTime = ((now.Hour() * 1000) - 6000 + (now.Minute() * 16))
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
	Home          `json:"home"`
	Weather       WeatherStatus     `json:"weather"`
	CurrentTime   int64             `json:"current_time"` // In Unix time.
	MinecraftTime int               `json:"minecraft_time"`
	ForecastData  *darksky.Forecast `json:"-"`
}

// WeatherStatus represents a weather status state in the simulated
// environment.
type WeatherStatus struct {
	Type        WeatherType `json:"type"`
	OutdoorTemp float64     `json:"outdoor_temp"` // In Celcius.
	Humidity    float64     `json:"humidity"`     // In decimal, 0.5 = 50%.
	CloudCover  float64     `json:"cloud_cover"`
}

// Home represents a simulated home state.
type Home struct {
	MainDoorOpened bool    `json:"main_door_opened"` // Whether the main door is opened or not.
	Rooms          []Room  `json:"rooms"`
	PowerGenRate   float64 `json:"power_gen_rate"`
	PowerConRate   float64 `json:"power_con_rate"`
	SolarMaxPower  int     `json:"solar_max_power"` // Maximum solar panel generation capacity, in kWh.
	NetPower       float64 `json:"net_power"`
	BatteryStore   float64 `json:"battery_store"`
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
