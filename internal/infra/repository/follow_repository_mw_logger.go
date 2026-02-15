package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
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
	x.LogMw(ctx, fields, err)

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
	x.LogMw(ctx, fields, err)

	return err
}
