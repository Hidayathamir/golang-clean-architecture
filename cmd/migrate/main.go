package main

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/l"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	viperConfig := config.NewViper()
	l.SetupLogger(viperConfig)

	config.Migrate(viperConfig)
}
