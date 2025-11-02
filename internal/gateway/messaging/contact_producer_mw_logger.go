package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ ContactProducer = &ContactProducerMwLogger{}

type ContactProducerMwLogger struct {
	Next ContactProducer
}

func NewContactProducerMwLogger(next ContactProducer) *ContactProducerMwLogger {
	return &ContactProducerMwLogger{
		Next: next,
	}
}

func (p *ContactProducerMwLogger) Send(ctx context.Context, event *model.ContactEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.Send(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	logging.Log(ctx, fields, err)

	return err
}
