package messaging_test

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"
)

// TestConsumeEventSingle_PollErrorDoesNotSpin reproduces a broker outage: with
// an unreachable broker every PollRecords returns an error. Before the backoff
// fix the loop `continue`d with no delay, so it iterated (and logged) hundreds
// of thousands of times per second. This asserts the loop is paced: within the
// window it can only iterate a bounded number of times, and it stops promptly on
// context cancellation.
func TestConsumeEventSingle_PollErrorDoesNotSpin(t *testing.T) {
	cfg := config.NewConfig()
	cfg.Set(config.KafkaBootstrapServers, "127.0.0.1:1") // nothing listens here

	producer := provider.NewKafkaClientProducer(cfg)
	defer producer.Close()

	// Count every "error client poll fetch" log line the loop emits.
	var pollErrors int64
	prevOut := logkit.Logger.Out
	logkit.Logger.SetOutput(logWriter{onLine: func() {
		atomic.AddInt64(&pollErrors, 1)
	}})
	defer logkit.Logger.SetOutput(prevOut)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		defer close(done)
		messaging.ConsumeEventSingle(ctx, cfg, producer, "test-group-spin", topic.UserFollowed, func(_ context.Context, _ *kgo.Record) error {
			return nil
		})
	}()

	// Let the consumer run against the dead broker for 2s.
	time.Sleep(2 * time.Second)
	cancel()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		require.Fail(t, "consumer did not stop after context cancellation")
	}

	got := atomic.LoadInt64(&pollErrors)
	// With backoff each attempt is >=100ms and duplicate errors collapse, so a
	// 2s outage yields a single-digit number of lines. Without the fix this is
	// in the hundreds of thousands. A generous ceiling still proves the point.
	assert.Less(t, got, int64(100), "poll-error loop is not paced (spun %d times in 2s)", got)
}

// logWriter counts the consumer's poll-fetch error lines. Logrus emits one
// Write call per log entry, so each call is a single line.
type logWriter struct {
	onLine func()
}

func (w logWriter) Write(p []byte) (int, error) {
	if strings.Contains(string(p), "error client poll fetch") && w.onLine != nil {
		w.onLine()
	}
	return len(p), nil
}
