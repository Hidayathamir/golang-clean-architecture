package search

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ ImageSearch2 = &ImageSearch2MwLogger{}

type ImageSearch2MwLogger struct {
	Next ImageSearch2
}

func NewImageSearch2MwLogger(next ImageSearch2) *ImageSearch2MwLogger {
	return &ImageSearch2MwLogger{
		Next: next,
	}
}

func (i *ImageSearch2MwLogger) IndexImage(ctx context.Context, document *dto.ImageDocument) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := i.Next.IndexImage(ctx, document)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"document": document,
	}
	x.LogMw(ctx, fields, err)

	return err
}
