package l

import (
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logrushook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Logger = logrus.New()

func SetupLogger(viperConfig *viper.Viper) {
	logger := logrus.New()

	logger.SetReportCaller(true)
	logger.SetLevel(logrus.Level(viperConfig.GetInt32(configkey.LogLevel)))
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.AddHook(logrushook.NewOtelHook())

	Logger = logger
}

func SetLogger(logger *logrus.Logger) {
	if logger == nil {
		logger = logrus.New()
	}
	Logger = logger
}
