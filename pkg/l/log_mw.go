package l

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/caller"
	"github.com/sirupsen/logrus"
)

func LogMw(ctx context.Context, fields logrus.Fields, err error) {
	level, errMsg := GetLevelAndErrMsg(err)

	fileLine, funcName := caller.Info(caller.WithSkip(1))

	Logger.WithContext(ctx).WithFields(logrus.Fields{
		"fields": LimitJSON(fields),
		"err":    errMsg,
		"source": fileLine,
	}).Log(level, funcName)
}

func GetLevelAndErrMsg(err error) (logrus.Level, string) {
	level := logrus.InfoLevel
	errMsg := ""
	if err != nil {
		level = logrus.ErrorLevel
		errMsg = err.Error()
	}
	return level, errMsg
}

var limitChar = 10000

func LimitJSON(v any) any {
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
