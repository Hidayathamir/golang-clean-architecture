package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ LikeRepository = &LikeRepositoryMwLogger{}

type LikeRepositoryMwLogger struct {
	Next LikeRepository
}

func NewLikeRepositoryMwLogger(next LikeRepository) *LikeRepositoryMwLogger {
	return &LikeRepositoryMwLogger{
		Next: next,
	}
}

func (r *LikeRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, like *entity.Like) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.Create(ctx, db, like)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"like": like,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *LikeRepositoryMwLogger) FindByImageID(ctx context.Context, db *gorm.DB, likeList *entity.LikeList, imageID int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindByImageID(ctx, db, likeList, imageID)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"likeList": likeList,
		"imageID":  imageID,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}
