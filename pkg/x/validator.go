package x

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func SetupValidator(cfg *config.Config) {
	Validate = validator.New()
}
