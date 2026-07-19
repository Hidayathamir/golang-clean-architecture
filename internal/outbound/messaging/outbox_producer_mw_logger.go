package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ OutboxProducer = &OutboxProducerMwLogger{}

type OutboxProducerMwLogger struct {
	Next OutboxProducer
}

func NewOutboxProducerMwLogger(next OutboxProducer) *OutboxProducerMwLogger {
	return &OutboxProducerMwLogger{
		Next: next,
	}
}

func (p *OutboxProducerMwLogger) ProducePending(ctx context.Context) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.ProducePending(ctx)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{}
	logkit.LogMw(ctx, fields, err)

	return err
}
