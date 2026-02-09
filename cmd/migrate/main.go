package main

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/migrate"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.NewConfig()
	x.SetupAll(cfg)

	db := provider.NewDatabase(cfg)

	migrate.Migrate(cfg, db)
}
