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

func (r *UserRepositoryMwLogger) CountByID(ctx context.Context, db *gorm.DB, id string) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	total, err := r.Next.CountByID(ctx, db, id)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":    id,
		"total": total,
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

func (r *UserRepositoryMwLogger) FindByID(ctx context.Context, db *gorm.DB, entity *entity.User, id string) error {
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
