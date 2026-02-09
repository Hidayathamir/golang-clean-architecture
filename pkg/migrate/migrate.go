package migrate

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func Migrate(cfg *config.Config, db *gorm.DB) {
	sqlDB, err := db.DB()
	x.PanicIfErr(err)

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	x.PanicIfErr(err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+cfg.GetDatabaseMigrations(),
		"postgres",
		driver,
	)
	x.PanicIfErr(err)

	err = m.Up()
	x.PanicIfErr(err)
}
