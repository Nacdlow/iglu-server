package sim

// ChangeTimeSleepForm represents a form to change the simulation's tick speed.
type ChangeTimeSleepForm struct {
	TickSleep int `form:"tick_sleep"`
}

// OverrideWeatherForm represents a form to change the simulation's weather.
type OverrideWeatherForm struct {
	OutdoorTemp float64 `form:"outdoor_temp"`
	Humidity    float64 `form:"humidity"`
	CloudCover  float64 `form:"cloud_cover"`
}
