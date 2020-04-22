package weather

import (
	"github.com/adlio/darksky"
	"github.com/Nacdlow/iglu-server/modules/settings"
)

// GetWeather loads the forecast from Darksky's API.
func GetWeather() (*darksky.Forecast, error) {
	lat := settings.Config.GetString("Location.Lat")
	lon := settings.Config.GetString("Location.Lon")
	client := darksky.NewClient(settings.Config.GetString("DarkskyAPIKey"))
	f, err := client.GetForecast(lat, lon,
		darksky.Arguments{"units": "si", "extend": "hourly"})
	if err != nil {
		return nil, err
	}
	return f, nil
}
