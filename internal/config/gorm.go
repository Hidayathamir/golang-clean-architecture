package config

import (
	"fmt"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/l"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viperConfig *viper.Viper) *gorm.DB {
	username := viperConfig.GetString(configkey.DatabaseUsername)
	password := viperConfig.GetString(configkey.DatabasePassword)
	host := viperConfig.GetString(configkey.DatabaseHost)
	port := viperConfig.GetInt(configkey.DatabasePort)
	database := viperConfig.GetString(configkey.DatabaseName)
	idleConnection := viperConfig.GetInt(configkey.DatabasePoolIdle)
	maxConnection := viperConfig.GetInt(configkey.DatabasePoolMax)
	maxLifeTimeConnection := viperConfig.GetInt(configkey.DatabasePoolLifetime)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC", host, port, username, password, database)

	gormLogger := logger.New(&logrusWriter{Logger: l.Logger}, logger.Config{
		SlowThreshold:             time.Second * 5,
		Colorful:                  false,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
		LogLevel:                  logger.Info,
	})

	const maxAttempts = 5

	var db *gorm.DB
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
		if err == nil {
			break
		}
		l.Logger.Warnf("database connection attempt %d/%d failed: %v", attempt, maxAttempts, err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		l.Logger.Panicf("failed to connect database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		l.Logger.Panicf("failed to connect database: %v", err)
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if pingErr := connection.Ping(); pingErr == nil {
			break
		} else if attempt == maxAttempts {
			l.Logger.Panicf("failed to connect database: %v", pingErr)
		} else {
			l.Logger.Warnf("database ping attempt %d/%d failed: %v", attempt, maxAttempts, pingErr)
			time.Sleep(1 * time.Second)
		}
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}
