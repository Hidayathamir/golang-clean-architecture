package config

import (
	"context"
	"fmt"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
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

	const maxAttempts = 5

	var db *gorm.DB
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:         &gormLogger{},
			TranslateError: true,
		})
		if err == nil {
			break
		}
		x.Logger.Warnf("database connection attempt %d/%d failed: %v", attempt, maxAttempts, err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		x.Logger.Panicf("failed to connect database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		x.Logger.Panicf("failed to connect database: %v", err)
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if pingErr := connection.Ping(); pingErr == nil {
			break
		} else if attempt == maxAttempts {
			x.Logger.Panicf("failed to connect database: %v", pingErr)
		} else {
			x.Logger.Warnf("database ping attempt %d/%d failed: %v", attempt, maxAttempts, pingErr)
			time.Sleep(1 * time.Second)
		}
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

type gormLogger struct{}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, msg string, data ...any) {
	x.Logger.WithContext(ctx).Infof(msg, data...)
}

func (l *gormLogger) Warn(ctx context.Context, msg string, data ...any) {
	x.Logger.WithContext(ctx).Warnf(msg, data...)
}

func (l *gormLogger) Error(ctx context.Context, msg string, data ...any) {
	x.Logger.WithContext(ctx).Errorf(msg, data...)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	entry := x.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"elapsed": time.Since(begin),
		"rows":    rows,
	})

	if err != nil {
		entry.Error(sql)
	} else {
		entry.Info(sql)
	}
}
