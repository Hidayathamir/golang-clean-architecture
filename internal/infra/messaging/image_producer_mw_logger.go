package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ ImageProducer = &ImageProducerMwLogger{}

type ImageProducerMwLogger struct {
	Next ImageProducer
}

func NewImageProducerMwLogger(next ImageProducer) *ImageProducerMwLogger {
	return &ImageProducerMwLogger{
		Next: next,
	}
}

func (p *ImageProducerMwLogger) SendImageUploaded(ctx context.Context, event *dto.ImageUploadedEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.SendImageUploaded(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (p *ImageProducerMwLogger) SendImageLiked(ctx context.Context, event *dto.ImageLikedEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.SendImageLiked(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (p *ImageProducerMwLogger) SendImageCommented(ctx context.Context, event *dto.ImageCommentedEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.SendImageCommented(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	x.LogMw(ctx, fields, err)

	return err
}
