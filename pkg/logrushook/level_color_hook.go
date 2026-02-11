package logrushook

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var _ logrus.Hook = &LevelColorHook{}

type LevelColorHook struct{}

func NewLevelColorHook() *LevelColorHook {
	return &LevelColorHook{}
}

func (h *LevelColorHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LevelColorHook) Fire(entry *logrus.Entry) error {
	levelColor := ""

	switch entry.Level {
	case logrus.PanicLevel:
		levelColor = "⚫"
	case logrus.FatalLevel:
		levelColor = "⚫"
	case logrus.ErrorLevel:
		levelColor = "🔴"
	case logrus.WarnLevel:
		levelColor = "🟠"
	case logrus.InfoLevel:
		levelColor = "🔵"
	case logrus.DebugLevel:
		levelColor = "🟢"
	case logrus.TraceLevel:
		levelColor = "🟣"
	}

	entry.Message = fmt.Sprintf("[%s] %s", levelColor, entry.Message)

	return nil
}
