package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ OutboxPublisher = &OutboxPublisherMwLogger{}

type OutboxPublisherMwLogger struct {
	Next OutboxPublisher
}

func NewOutboxPublisherMwLogger(next OutboxPublisher) *OutboxPublisherMwLogger {
	return &OutboxPublisherMwLogger{
		Next: next,
	}
}

func (p *OutboxPublisherMwLogger) PublishPending(ctx context.Context) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.PublishPending(ctx)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{}
	logkit.LogMw(ctx, fields, err)

	return err
}
