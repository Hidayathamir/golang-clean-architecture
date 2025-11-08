package x

import (
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logrushook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Logger = logrus.New()

func SetupLogger(viperConfig *viper.Viper) {
	logger := logrus.New()

	lvl, err := logrus.ParseLevel(viperConfig.GetString(configkey.LogLevel))
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logger.SetReportCaller(true)
	logger.SetLevel(lvl)
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
