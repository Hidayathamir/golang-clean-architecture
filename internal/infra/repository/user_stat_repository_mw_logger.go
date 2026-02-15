package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ UserStatRepository = &UserStatRepositoryMwLogger{}

type UserStatRepositoryMwLogger struct {
	Next UserStatRepository
}

func NewUserStatRepositoryMwLogger(next UserStatRepository) *UserStatRepositoryMwLogger {
	return &UserStatRepositoryMwLogger{
		Next: next,
	}
}

func (r *UserStatRepositoryMwLogger) IncrementFollowerCountAndFollowingCountByID(ctx context.Context, db *gorm.DB, id int64, followerCount int, followingCount int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.IncrementFollowerCountAndFollowingCountByID(ctx, db, id, followerCount, followingCount)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":             id,
		"followerCount":  followerCount,
		"followingCount": followingCount,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserStatRepositoryMwLogger) IncrementFollowerCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.IncrementFollowerCountByID(ctx, db, id, count)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":    id,
		"count": count,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserStatRepositoryMwLogger) IncrementFollowingCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.IncrementFollowingCountByID(ctx, db, id, count)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":    id,
		"count": count,
	}
	x.LogMw(ctx, fields, err)

	return err
}
