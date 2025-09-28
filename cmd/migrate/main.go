package main

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)

	config.Migrate(viperConfig, log)
}
