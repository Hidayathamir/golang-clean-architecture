package retrykit

import (
	"context"
	"errors"
	"net"
	"syscall"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/avast/retry-go/v5"
)

// Per-layer retry tuning. DB and cache are local and fast to retry; messaging
// (Kafka) and search (Elasticsearch) sit behind a broker and need more backoff;
// storage (S3) round-trips are slow, so retry less aggressively.
var (
	dbRetryOption = []retry.Option{
		retry.Attempts(3),
		retry.Delay(500 * time.Millisecond),
		retry.RetryIf(retryIf),
	}
	cacheRetryOption = []retry.Option{
		retry.Attempts(3),
		retry.Delay(50 * time.Millisecond),
		retry.RetryIf(retryIf),
	}
	messagingRetryOption = []retry.Option{
		retry.Attempts(5),
		retry.Delay(1 * time.Second),
		retry.MaxDelay(10 * time.Second),
		retry.DelayType(retry.BackOffDelay),
		retry.RetryIf(retryIf),
	}
	searchRetryOption = []retry.Option{
		retry.Attempts(4),
		retry.Delay(500 * time.Millisecond),
		retry.MaxDelay(5 * time.Second),
		retry.DelayType(retry.BackOffDelay),
		retry.RetryIf(retryIf),
	}
	storageRetryOption = []retry.Option{
		retry.Attempts(3),
		retry.Delay(1 * time.Second),
		retry.MaxDelay(5 * time.Second),
		retry.DelayType(retry.BackOffDelay),
		retry.RetryIf(retryIf),
	}
)

func DBRetry(ctx context.Context, operation func() error) error {
	return retry.New(withContext(ctx, dbRetryOption, "DBRetry")...).
		Do(func() error { return operation() })
}

func DBRetryWithData[T any](ctx context.Context, operation func() (T, error)) (T, error) {
	return retry.NewWithData[T](withContext(ctx, dbRetryOption, "DBRetryWithData")...).
		Do(func() (T, error) { return operation() })
}

func CacheRetry(ctx context.Context, operation func() error) error {
	return retry.New(withContext(ctx, cacheRetryOption, "CacheRetry")...).
		Do(func() error { return operation() })
}

func CacheRetryWithData[T any](ctx context.Context, operation func() (T, error)) (T, error) {
	return retry.NewWithData[T](withContext(ctx, cacheRetryOption, "CacheRetryWithData")...).
		Do(func() (T, error) { return operation() })
}

func MessagingRetry(ctx context.Context, operation func() error) error {
	return retry.New(withContext(ctx, messagingRetryOption, "MessagingRetry")...).
		Do(func() error { return operation() })
}

func MessagingRetryWithData[T any](ctx context.Context, operation func() (T, error)) (T, error) {
	return retry.NewWithData[T](withContext(ctx, messagingRetryOption, "MessagingRetryWithData")...).
		Do(func() (T, error) { return operation() })
}

func SearchRetry(ctx context.Context, operation func() error) error {
	return retry.New(withContext(ctx, searchRetryOption, "SearchRetry")...).
		Do(func() error { return operation() })
}

func SearchRetryWithData[T any](ctx context.Context, operation func() (T, error)) (T, error) {
	return retry.NewWithData[T](withContext(ctx, searchRetryOption, "SearchRetryWithData")...).
		Do(func() (T, error) { return operation() })
}

func StorageRetry(ctx context.Context, operation func() error) error {
	return retry.New(withContext(ctx, storageRetryOption, "StorageRetry")...).
		Do(func() error { return operation() })
}

func StorageRetryWithData[T any](ctx context.Context, operation func() (T, error)) (T, error) {
	return retry.NewWithData[T](withContext(ctx, storageRetryOption, "StorageRetryWithData")...).
		Do(func() (T, error) { return operation() })
}

// withContext appends the shared per-attempt options to a layer's base options.
func withContext(ctx context.Context, base []retry.Option, name string) []retry.Option {
	return append(
		base,
		retry.Context(ctx),
		retry.OnRetry(func(n uint, err error) {
			logkit.Logger.WithError(err).WithField("attempt", n+1).Warn(name)
		}),
	)
}

// -------------------------------------------------------------------------- //

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
