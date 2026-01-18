package x

import "context"

func LogIfErr(ctx context.Context, err error) {
	if err != nil {
		Logger.WithContext(ctx).WithError(err).Info()
	}
}
