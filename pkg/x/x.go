package x

import "github.com/spf13/viper"

func SetupAll(viperConfig *viper.Viper) {
	SetupLogger(viperConfig)
	SetupValidator(viperConfig)
}
