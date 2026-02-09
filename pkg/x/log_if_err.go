package x

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/caller"
	"github.com/sirupsen/logrus"
)

func LogIfErr(err error) {
	LogIfErrContext(context.Background(), err, 1)
}

func LogIfErrContext(ctx context.Context, err error, skips ...int) {
	if err != nil {
		skip := 1
		if len(skips) > 0 {
			skip += skips[0]
		}
		fileLine, funcName := caller.Info(caller.WithSkip(skip))
		Logger.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
			"callerSource":   fileLine,
			"callerFuncName": funcName,
		}).Warn()
	}
}

func LogIfErrForDefer(f func() error) {
	LogIfErrForDeferContext(context.Background(), f, 1)
}

func LogIfErrForDeferContext(ctx context.Context, f func() error, skips ...int) {
	err := f()
	if err != nil {
		skip := 1
		if len(skips) > 0 {
			skip += skips[0]
		}
		fileLine, funcName := caller.Info(caller.WithSkip(skip))
		Logger.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
			"callerSource":   fileLine,
			"callerFuncName": funcName,
		}).Warn()
	}
}
