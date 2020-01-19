package settings

import (
	"github.com/spf13/viper"
	"gitlab.com/skilstak/code/go/uniq"
)

var (
	Config = viper.New()
)

func LoadConfig() {
	Config.SetConfigName("config")
	Config.AddConfigPath("/etc/nacdlow/")
	Config.AddConfigPath("$HOME/.config/nacdlow")
	Config.AddConfigPath(".")

	Config.SetDefault("Port", "8080")
	Config.SetDefault("HouseName", "My House")
	Config.SetDefault("Address", "Heriot-Watt University, Edinburgh, Scotland. EH14 4AS")
	Config.SetDefault("Engineer.Name", "Nacdlow Support")
	Config.SetDefault("Engineer.Phone", "074 123 4567")
	Config.SetDefault("Engineer.Email", "support@nacdlow.com")
	if !Config.IsSet("CustomerID") {
		Config.SetDefault("CustomerID", uniq.UUID())
	}
	//Config.SetDefault("","")

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

}
