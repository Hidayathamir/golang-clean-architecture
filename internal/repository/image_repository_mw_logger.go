package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ ImageRepository = &ImageRepositoryMwLogger{}

type ImageRepositoryMwLogger struct {
	Next ImageRepository
}

func NewImageRepositoryMwLogger(next ImageRepository) *ImageRepositoryMwLogger {
	return &ImageRepositoryMwLogger{
		Next: next,
	}
}

func (r *ImageRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, image *entity.Image) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Create(ctx, db, image)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"image": image,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *ImageRepositoryMwLogger) FindByID(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByID(ctx, db, image, id)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"image": image,
		"id":    id,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *ImageRepositoryMwLogger) IncrementCommentCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.IncrementCommentCountByID(ctx, db, id, count)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":    id,
		"count": count,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *ImageRepositoryMwLogger) IncrementLikeCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.IncrementLikeCountByID(ctx, db, id, count)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":    id,
		"count": count,
	}
	x.LogMw(ctx, fields, err)

	return err
}
