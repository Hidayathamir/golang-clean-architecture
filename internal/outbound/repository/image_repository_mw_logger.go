package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
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

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.Create(ctx, db, image)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"image": image,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *ImageRepositoryMwLogger) FindByID(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindByID(ctx, db, image, id)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"image": image,
		"id":    id,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *ImageRepositoryMwLogger) FindFeedByFollowerID(ctx context.Context, db *gorm.DB, imageList *entity.ImageList, followerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindFeedByFollowerID(ctx, db, imageList, followerID, cursorCreatedAt, cursorID, limit)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"image_list": imageList,
		"follower":   followerID,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *ImageRepositoryMwLogger) FindByUserID(ctx context.Context, db *gorm.DB, imageList *entity.ImageList, userID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindByUserID(ctx, db, imageList, userID, cursorCreatedAt, cursorID, limit)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"image_list": imageList,
		"user":       userID,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *ImageRepositoryMwLogger) IncrementCommentCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.IncrementCommentCountByID(ctx, db, id, count)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":    id,
		"count": count,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *ImageRepositoryMwLogger) IncrementLikeCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.IncrementLikeCountByID(ctx, db, id, count)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":    id,
		"count": count,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}
