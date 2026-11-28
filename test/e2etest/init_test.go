package e2etest

import (
	"os"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/migrate"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var db *gorm.DB

var cfg *config.Config

var redisClient *redis.Client

// TestMain is the entry point for all tests in this package.
// It sets up global dependencies (logger, validator, Fiber app, DB, Kafka producer),
// starts a PostgreSQL container, runs database migrations, bootstraps the application,
// executes test, and finally terminates the container before exiting.
func TestMain(m *testing.M) {
	cfg = config.NewConfig()
	logkit.SetupLogger(cfg)
	validatorkit.SetupValidator(cfg)
	db = provider.NewDatabase(cfg)
	// The app servers share this database, so instead of dropping the schema
	// (which blocks on their open transactions) each test clears its data via
	// ClearAll. Migrations here are idempotent and apply 0 when already current.
	cfg.SetDatabaseMigrations("../../db/migrations")
	migrate.Migrate(cfg, db)

	redisClient = provider.NewRedisClient(cfg)

	code := m.Run()

	os.Exit(code)
}
