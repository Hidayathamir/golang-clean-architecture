package migrate

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm"
)

func Migrate(cfg *config.Config, db *gorm.DB) {
	sqlDB, err := db.DB()
	x.PanicIfErr(err)

	migrations := &migrate.FileMigrationSource{
		Dir: cfg.GetDatabaseMigrations(),
	}

	n, err := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	x.PanicIfErr(err)

	x.Logger.Infof("Applied %d migrations", n)
}
