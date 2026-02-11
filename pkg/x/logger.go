package x

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logrushook"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func SetupLogger(cfg *config.Config) {
	logger := logrus.New()

	lvl, err := logrus.ParseLevel(cfg.GetLogLevel())
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logger.SetReportCaller(true)
	logger.SetLevel(lvl)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.AddHook(logrushook.NewLevelColorHook())
	logger.AddHook(logrushook.NewOtelHook())

	Logger = logger
}

func SetLogger(logger *logrus.Logger) {
	if logger == nil {
		logger = logrus.New()
	}
	Logger = logger
}
