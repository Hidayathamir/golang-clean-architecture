package migrate

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm"
)

func Migrate(cfg *config.Config, db *gorm.DB) {
	sqlDB, err := db.DB()
	errkit.PanicIfErr(err)

	migrations := &migrate.FileMigrationSource{
		Dir: cfg.GetDatabaseMigrations(),
	}

	n, err := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	errkit.PanicIfErr(err)

	logkit.Logger.Infof("Applied %d migrations", n)
}
