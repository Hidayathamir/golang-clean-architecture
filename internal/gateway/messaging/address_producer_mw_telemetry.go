package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ AddressProducer = &AddressProducerMwTelemetry{}

type AddressProducerMwTelemetry struct {
	Next AddressProducer
}

func NewAddressProducerMwTelemetry(next AddressProducer) *AddressProducerMwTelemetry {
	return &AddressProducerMwTelemetry{
		Next: next,
	}
}

func (p *AddressProducerMwTelemetry) Send(ctx context.Context, event *model.AddressEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.Send(ctx, event)
	telemetry.RecordError(span, err)

	return err
}
