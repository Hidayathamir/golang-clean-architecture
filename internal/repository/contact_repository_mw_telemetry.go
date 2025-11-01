package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"gorm.io/gorm"
)

var _ ContactRepository = &ContactRepositoryMwTelemetry{}

type ContactRepositoryMwTelemetry struct {
	Next ContactRepository
}

func NewContactRepositoryMwTelemetry(next ContactRepository) *ContactRepositoryMwTelemetry {
	return &ContactRepositoryMwTelemetry{
		Next: next,
	}
}

func (r *ContactRepositoryMwTelemetry) FindByIDAndUserID(ctx context.Context, db *gorm.DB, contact *entity.Contact, id string, userID string) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByIDAndUserID(ctx, db, contact, id, userID)
	telemetry.RecordError(span, err)

	return err
}

func (r *ContactRepositoryMwTelemetry) Search(ctx context.Context, db *gorm.DB, req *model.SearchContactRequest) (entity.ContactList, int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	list, total, err := r.Next.Search(ctx, db, req)
	telemetry.RecordError(span, err)

	return list, total, err
}

func (r *ContactRepositoryMwTelemetry) Create(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Create(ctx, db, entity)
	telemetry.RecordError(span, err)

	return err
}

func (r *ContactRepositoryMwTelemetry) Update(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Update(ctx, db, entity)
	telemetry.RecordError(span, err)

	return err
}

func (r *ContactRepositoryMwTelemetry) Delete(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Delete(ctx, db, entity)
	telemetry.RecordError(span, err)

	return err
}
