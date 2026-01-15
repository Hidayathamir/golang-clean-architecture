package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ ImageUploadedProducer = &ImageUploadedProducerMwLogger{}

type ImageUploadedProducerMwLogger struct {
	Next ImageUploadedProducer
}

func NewImageUploadedProducerMwLogger(next ImageUploadedProducer) *ImageUploadedProducerMwLogger {
	return &ImageUploadedProducerMwLogger{
		Next: next,
	}
}

func (p *ImageUploadedProducerMwLogger) Send(ctx context.Context, event *model.ImageUploadedEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.Send(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	x.LogMw(ctx, fields, err)

	return err
}
