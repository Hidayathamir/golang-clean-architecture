package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ UserRepository = &UserRepositoryMwLogger{}

type UserRepositoryMwLogger struct {
	Next UserRepository
}

func NewUserRepositoryMwLogger(next UserRepository) *UserRepositoryMwLogger {
	return &UserRepositoryMwLogger{
		Next: next,
	}
}

func (r *UserRepositoryMwLogger) CountByUsername(ctx context.Context, db *gorm.DB, username string) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	total, err := r.Next.CountByUsername(ctx, db, username)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"username": username,
		"total":    total,
	}
	x.LogMw(ctx, fields, err)

	return total, err
}

func (r *UserRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, user *entity.User) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Create(ctx, db, user)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"user": user,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) FindByID(ctx context.Context, db *gorm.DB, user *entity.User, id int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByID(ctx, db, user, id)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":   id,
		"user": user,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) FindByUsername(ctx context.Context, db *gorm.DB, user *entity.User, username string) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByUsername(ctx, db, user, username)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"username": username,
		"user":     user,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, user *entity.User) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Update(ctx, db, user)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"user": user,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) IncrementFollowerCountAndFollowingCountByID(ctx context.Context, db *gorm.DB, id int64, followerCount int, followingCount int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.IncrementFollowerCountAndFollowingCountByID(ctx, db, id, followerCount, followingCount)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":             id,
		"followerCount":  followerCount,
		"followingCount": followingCount,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) IncrementFollowerCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.IncrementFollowerCountByID(ctx, db, id, count)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":    id,
		"count": count,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) IncrementFollowingCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.IncrementFollowingCountByID(ctx, db, id, count)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":    id,
		"count": count,
	}
	x.LogMw(ctx, fields, err)

	return err
}
