package x

import "github.com/Hidayathamir/golang-clean-architecture/internal/config"

func SetupAll(cfg *config.Config) {
	SetupLogger(cfg)
	SetupValidator(cfg)
}
