package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ ContactProducer = &ContactProducerMwTelemetry{}

type ContactProducerMwTelemetry struct {
	Next ContactProducer
}

func NewContactProducerMwTelemetry(next ContactProducer) *ContactProducerMwTelemetry {
	return &ContactProducerMwTelemetry{
		Next: next,
	}
}

func (p *ContactProducerMwTelemetry) Send(ctx context.Context, event *model.ContactEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.Send(ctx, event)
	telemetry.RecordError(span, err)

	return err
}
