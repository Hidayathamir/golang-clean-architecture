package config

import (
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Migrate(viperConfig *viper.Viper, log *logrus.Logger) {
	db := NewDatabase(viperConfig, log)

	sqlDB, err := db.DB()
	panicIfErr(err)

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	panicIfErr(err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+viperConfig.GetString(configkey.DatabaseMigrations),
		"postgres",
		driver,
	)
	panicIfErr(err)

	err = m.Up()
	panicIfErr(err)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
