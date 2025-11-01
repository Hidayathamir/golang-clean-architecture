package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ UserProducer = &UserProducerMwTelemetry{}

type UserProducerMwTelemetry struct {
	Next UserProducer
}

func NewUserProducerMwTelemetry(next UserProducer) *UserProducerMwTelemetry {
	return &UserProducerMwTelemetry{
		Next: next,
	}
}

func (p *UserProducerMwTelemetry) Send(ctx context.Context, event *model.UserEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.Send(ctx, event)
	telemetry.RecordError(span, err)

	return err
}
