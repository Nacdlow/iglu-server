package sim

type ChangeTimeSleepForm struct {
	TickSleep int `form:"tick_sleep" binding:"Required"`
}

type OverrideWeatherForm struct {
	OutdoorTemp float64 `form:"outdoor_temp"`
	Humidity    float64 `form:"humidity"`
	CloudCover  float64 `form:"cloud_cover"`
}
