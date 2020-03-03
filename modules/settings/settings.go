package settings

import (
	"github.com/spf13/viper"
	"gitlab.com/skilstak/code/go/uniq"
)

var (
	// Config is the viper configuration file.
	Config = viper.New()
)

// LoadConfig loads the configuration file and sets the default settings.
func LoadConfig() {
	Config.SetConfigName("config")
	Config.AddConfigPath("/etc/nacdlow/")
	Config.AddConfigPath("$HOME/.config/nacdlow")
	Config.AddConfigPath(".")

	// Fill configuration with default
	Config.SetDefault("Port", "8080")
	Config.SetDefault("HouseName", "My House")
	Config.SetDefault("Address", "Heriot-Watt University, Edinburgh, Scotland. EH14 4AS")
	Config.SetDefault("Engineer.Name", "Nacdlow Support")
	Config.SetDefault("Engineer.Phone", "0131 496 0143")
	Config.SetDefault("Engineer.Email", "support@nacdlow.com")
	if !Config.IsSet("CustomerID") {
		Config.SetDefault("CustomerID", uniq.UUID())
	}
	Config.SetDefault("Timezone", "Europe/London")
	Config.SetDefault("DarkskyAPIKey", "APIKEYHERE")
	Config.SetDefault("Location.Lat", "25.371679")
	Config.SetDefault("Location.Lon", "55.511716")
	Config.SetDefault("Simulation.SolarCapacityKW", 45)
	Config.SetDefault("Simulation.BatteryCapacityKWH", 135)
	Config.SetDefault("Marketplace.RepositoryURL", "https://market.nacdlow.com/repo")

	// Read configuration
	if err := Config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := Config.WriteConfigAs("config.toml")
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	// We write to make sure that the default configuration values are stored
	err := Config.WriteConfig()
	if err != nil {
		panic(err)
	}
}
