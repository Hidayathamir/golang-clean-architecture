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

var _ FollowRepository = &FollowRepositoryMwLogger{}

type FollowRepositoryMwLogger struct {
	Next FollowRepository
}

func NewFollowRepositoryMwLogger(next FollowRepository) *FollowRepositoryMwLogger {
	return &FollowRepositoryMwLogger{
		Next: next,
	}
}

func (r *FollowRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, follow *entity.Follow) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.Create(ctx, db, follow)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"follow": follow,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *FollowRepositoryMwLogger) FindByFollowingID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followingID int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindByFollowingID(ctx, db, followList, followingID)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"followList":  followList,
		"followingID": followingID,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *FollowRepositoryMwLogger) DeleteByFollowerIDAndFollowingID(ctx context.Context, db *gorm.DB, followerID int64, followingID int64) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	deleted, err := retrykit.DBRetryWithData(ctx, func() (int64, error) {
		return r.Next.DeleteByFollowerIDAndFollowingID(ctx, db, followerID, followingID)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"followerID":  followerID,
		"followingID": followingID,
		"deleted":     deleted,
	}
	logkit.LogMw(ctx, fields, err)

	return deleted, err
}

func (r *FollowRepositoryMwLogger) FindFollowersByFollowingID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followingID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindFollowersByFollowingID(ctx, db, followList, followingID, cursorCreatedAt, cursorID, limit)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"followList":  followList,
		"followingID": followingID,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *FollowRepositoryMwLogger) FindFollowingsByFollowerID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindFollowingsByFollowerID(ctx, db, followList, followerID, cursorCreatedAt, cursorID, limit)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"followList": followList,
		"followerID": followerID,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}
