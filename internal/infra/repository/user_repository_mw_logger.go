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

	total, err := retrykit.DBRetryWithData(ctx, func() (int64, error) {
		return r.Next.CountByUsername(ctx, db, username)
	})
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

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.Create(ctx, db, user)
	})
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

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindByID(ctx, db, user, id)
	})
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

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindByUsername(ctx, db, user, username)
	})
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

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.Update(ctx, db, user)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"user": user,
	}
	x.LogMw(ctx, fields, err)

	return err
}
