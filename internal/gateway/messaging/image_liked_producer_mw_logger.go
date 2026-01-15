package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ ImageLikedProducer = &ImageLikedProducerMwLogger{}

type ImageLikedProducerMwLogger struct {
	Next ImageLikedProducer
}

func NewImageLikedProducerMwLogger(next ImageLikedProducer) *ImageLikedProducerMwLogger {
	return &ImageLikedProducerMwLogger{
		Next: next,
	}
}

func (p *ImageLikedProducerMwLogger) Send(ctx context.Context, event *model.ImageLikedEvent) error {
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
