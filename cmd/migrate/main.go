package main

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/migrate"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func main() {
	cfg := config.NewConfig()
	logkit.SetupLogger(cfg)
	validatorkit.SetupValidator(cfg)

	db := provider.NewDatabase(cfg)

	migrate.Migrate(cfg, db)
}
