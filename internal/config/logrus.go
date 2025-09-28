package config

import (
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/helper"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logrushook"
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
