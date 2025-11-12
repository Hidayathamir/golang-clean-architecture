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

func (r *UserRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Create(ctx, db, entity)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"entity": entity,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) FindByID(ctx context.Context, db *gorm.DB, entity *entity.User, id int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByID(ctx, db, entity, id)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":     id,
		"entity": entity,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) FindByUsername(ctx context.Context, db *gorm.DB, entity *entity.User, username string) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByUsername(ctx, db, entity, username)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"username": username,
		"entity":   entity,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Update(ctx, db, entity)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"entity": entity,
	}
	x.LogMw(ctx, fields, err)

	return err
}
