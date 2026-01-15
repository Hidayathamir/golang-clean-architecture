package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ ImageCommentedProducer = &ImageCommentedProducerMwLogger{}

type ImageCommentedProducerMwLogger struct {
	Next ImageCommentedProducer
}

func NewImageCommentedProducerMwLogger(next ImageCommentedProducer) *ImageCommentedProducerMwLogger {
	return &ImageCommentedProducerMwLogger{
		Next: next,
	}
}

func (p *ImageCommentedProducerMwLogger) Send(ctx context.Context, event *model.ImageCommentedEvent) error {
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
