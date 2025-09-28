package config

import (
	"fmt"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString(configkey.DatabaseUsername)
	password := viper.GetString(configkey.DatabasePassword)
	host := viper.GetString(configkey.DatabaseHost)
	port := viper.GetInt(configkey.DatabasePort)
	database := viper.GetString(configkey.DatabaseName)
	idleConnection := viper.GetInt(configkey.DatabasePoolIdle)
	maxConnection := viper.GetInt(configkey.DatabasePoolMax)
	maxLifeTimeConnection := viper.GetInt(configkey.DatabasePoolLifetime)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Panicf("failed to connect database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Panicf("failed to connect database: %v", err)
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
