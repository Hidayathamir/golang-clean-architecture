package x

import "context"

func LogIfErr(err error) {
	LogIfErrContext(context.Background(), err)
}

func LogIfErrContext(ctx context.Context, err error) {
	if err != nil {
		Logger.WithContext(ctx).WithError(err).Warn()
	}
}

func LogIfErrForDefer(f func() error) {
	LogIfErrForDeferContext(context.Background(), f)
}

func LogIfErrForDeferContext(ctx context.Context, f func() error) {
	err := f()
	if err != nil {
		Logger.WithContext(ctx).WithError(err).Warn()
	}
}
