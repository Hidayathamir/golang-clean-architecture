package retrykit

import (
	"context"
	"errors"
	"net"
	"syscall"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/avast/retry-go/v5"
)

func DBRetry(ctx context.Context, operation func() error) error {
	opts := append(
		retryOption,
		retry.Context(ctx),
		retry.OnRetry(func(n uint, err error) {
			x.Logger.WithError(err).WithField("attempt", n+1).Warn("DBRetry")
		}),
	)
	return retry.New(opts...).Do(func() error { return operation() })
}

func DBRetryWithData[T any](ctx context.Context, operation func() (T, error)) (T, error) {
	opts := append(
		retryOption,
		retry.Context(ctx),
		retry.OnRetry(func(n uint, err error) {
			x.Logger.WithError(err).WithField("attempt", n+1).Warn("DBRetryWithData")
		}),
	)
	return retry.NewWithData[T](opts...).Do(func() (T, error) { return operation() })
}

// -------------------------------------------------------------------------- //

var retryOption = []retry.Option{
	retry.Attempts(3),
	retry.Delay(500 * time.Millisecond),
	retry.RetryIf(retryIf),
}

func retryIf(err error) bool {
	if err == nil {
		return false
	}

	if isNetworkConnectionError(err) {
		return true
	}

	if isContextCancel(err) {
		return false
	}

	return false
}

func isNetworkConnectionError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && (netErr.Timeout()) {
		return true
	}
	if errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ETIMEDOUT) || errors.Is(err, syscall.ECONNREFUSED) {
		return true
	}
	return false
}

func isContextCancel(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
