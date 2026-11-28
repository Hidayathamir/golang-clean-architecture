package search

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
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

	err := retrykit.SearchRetry(ctx, func() error {
		return i.Next.IndexImage(ctx, document)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"document": document,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (i *ImageSearchMwLogger) SearchImages(ctx context.Context, imageList *entity.ImageList, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.SearchRetry(ctx, func() error {
		return i.Next.SearchImages(ctx, imageList, cursorCreatedAt, cursorID, limit)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"imageList": imageList,
		"limit":     limit,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}
