package x

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var Validate = validator.New()

func SetupValidator(viperConfig *viper.Viper) {
	Validate = validator.New()
}
