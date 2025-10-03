package helper

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/caller"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func SetLogger(l *logrus.Logger) { logger = l }

func Log(ctx context.Context, fields logrus.Fields, err error) {
	level, errMsg := getLevelAndErrMsg(err)

	logger.WithContext(ctx).WithFields(logrus.Fields{
		"fields": limitJSON(fields),
		"err":    errMsg,
		"source": caller.FileLine(caller.WithSkip(1)),
	}).Log(level, caller.FuncName(caller.WithSkip(1)))
}

func getLevelAndErrMsg(err error) (logrus.Level, string) {
	level := logrus.InfoLevel
	errMsg := ""
	if err != nil {
		level = logrus.ErrorLevel
		errMsg = err.Error()
	}
	return level, errMsg
}

var limitChar = 10000

func limitJSON(v any) any {
	jsonByte, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	jsonStr := string(jsonByte)
	if len(jsonStr) > limitChar {
		jsonStr = jsonStr[:limitChar] + "..."
		return jsonStr
	}
	return v
}
