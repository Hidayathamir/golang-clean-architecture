package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"gorm.io/gorm"
)

var _ UserRepository = &UserRepositoryMwTelemetry{}

type UserRepositoryMwTelemetry struct {
	Next UserRepository
}

func NewUserRepositoryMwTelemetry(next UserRepository) *UserRepositoryMwTelemetry {
	return &UserRepositoryMwTelemetry{
		Next: next,
	}
}

func (r *UserRepositoryMwTelemetry) FindByToken(ctx context.Context, db *gorm.DB, user *entity.User, token string) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByToken(ctx, db, user, token)
	telemetry.RecordError(span, err)

	return err
}

func (r *UserRepositoryMwTelemetry) Create(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Create(ctx, db, entity)
	telemetry.RecordError(span, err)

	return err
}

func (r *UserRepositoryMwTelemetry) Update(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Update(ctx, db, entity)
	telemetry.RecordError(span, err)

	return err
}

func (r *UserRepositoryMwTelemetry) CountByID(ctx context.Context, db *gorm.DB, id string) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	count, err := r.Next.CountByID(ctx, db, id)
	telemetry.RecordError(span, err)

	return count, err
}

func (r *UserRepositoryMwTelemetry) FindByID(ctx context.Context, db *gorm.DB, entity *entity.User, id string) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByID(ctx, db, entity, id)
	telemetry.RecordError(span, err)

	return err
}
