package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/l"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ AddressProducer = &AddressProducerMwLogger{}

type AddressProducerMwLogger struct {
	Next AddressProducer
}

func NewAddressProducerMwLogger(next AddressProducer) *AddressProducerMwLogger {
	return &AddressProducerMwLogger{
		Next: next,
	}
}

func (p *AddressProducerMwLogger) Send(ctx context.Context, event *model.AddressEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.Send(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	l.LogMw(ctx, fields, err)

	return err
}
