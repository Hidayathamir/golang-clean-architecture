package config

import (
	"golang-clean-architecture/pkg/constant/configkey"
	"golang-clean-architecture/pkg/helper"
	"golang-clean-architecture/pkg/logrushook"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogger(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(viper.GetInt32(configkey.LogLevel)))
	log.SetFormatter(&logrus.JSONFormatter{})

	log.AddHook(logrushook.NewTraceID())

	helper.SetLogger(log)
	return log
}
