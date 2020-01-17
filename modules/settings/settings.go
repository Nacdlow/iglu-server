package settings

import (
	"github.com/spf13/viper"
)

var (
	Config = viper.New()
)

func LoadConfig() {
	Config.SetConfigName("config")
	Config.AddConfigPath("/etc/nacdlow/")
	Config.AddConfigPath("$HOME/.config/nacdlow")
	Config.AddConfigPath(".")

	Config.SetDefault("HouseName", "My House")
	Config.SetDefault("Address", "Heriot-Watt University, Edinburgh, Scotland. EH14 4AS")
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
