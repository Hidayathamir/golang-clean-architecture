package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func NewConfig() *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath("./../../")
	v.AddConfigPath("./../")
	v.AddConfigPath("./")

	config := &Config{Viper: v}

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return config
}
