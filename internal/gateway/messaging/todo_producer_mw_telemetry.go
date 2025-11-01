package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ TodoProducer = &TodoProducerMwTelemetry{}

type TodoProducerMwTelemetry struct {
	Next TodoProducer
}

func NewTodoProducerMwTelemetry(next TodoProducer) *TodoProducerMwTelemetry {
	return &TodoProducerMwTelemetry{
		Next: next,
	}
}

func (p *TodoProducerMwTelemetry) Send(ctx context.Context, event *model.TodoCompletedEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.Send(ctx, event)
	telemetry.RecordError(span, err)

	return err
}
