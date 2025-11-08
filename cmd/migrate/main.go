package main

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	viperConfig := config.NewViper()
	x.SetupAll(viperConfig)

	config.Migrate(viperConfig)
}
