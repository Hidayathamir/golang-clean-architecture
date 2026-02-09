package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ NotifProducer = &NotifProducerMwLogger{}

type NotifProducerMwLogger struct {
	Next NotifProducer
}

func NewNotifProducerMwLogger(next NotifProducer) *NotifProducerMwLogger {
	return &NotifProducerMwLogger{
		Next: next,
	}
}

func (p *NotifProducerMwLogger) SendNotif(ctx context.Context, event *dto.NotifEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.SendNotif(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	x.LogMw(ctx, fields, err)

	return err
}
