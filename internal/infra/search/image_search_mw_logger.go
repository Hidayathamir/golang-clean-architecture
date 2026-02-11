package search

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ ImageSearch = &ImageSearchMwLogger{}

type ImageSearchMwLogger struct {
	Next ImageSearch
}

func NewImageSearchMwLogger(next ImageSearch) *ImageSearchMwLogger {
	return &ImageSearchMwLogger{
		Next: next,
	}
}

func (i *ImageSearchMwLogger) IndexImage(ctx context.Context, document *dto.ImageDocument) error {
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
