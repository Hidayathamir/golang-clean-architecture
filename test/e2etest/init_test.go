package e2etest

import (
	"os"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB

var viperConfig *viper.Viper

// TestMain is the entry point for all tests in this package.
// It sets up global dependencies (logger, validator, Fiber app, DB, Kafka producer),
// starts a PostgreSQL container, runs database migrations, bootstraps the application,
// executes test, and finally terminates the container before exiting.
func TestMain(m *testing.M) {
	viperConfig = config.NewViper()
	x.SetupAll(viperConfig)
	db = config.NewDatabase(viperConfig)
	db.Exec(`DROP SCHEMA public CASCADE`)
	db.Exec(`CREATE SCHEMA public`)
	viperConfig.Set(configkey.DatabaseMigrations, "../../db/migrations")
	config.Migrate(viperConfig)

	code := m.Run()

	os.Exit(code)
}
