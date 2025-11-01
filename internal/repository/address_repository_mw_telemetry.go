package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"gorm.io/gorm"
)

var _ AddressRepository = &AddressRepositoryMwTelemetry{}

type AddressRepositoryMwTelemetry struct {
	Next AddressRepository
}

func NewAddressRepositoryMwTelemetry(next AddressRepository) *AddressRepositoryMwTelemetry {
	return &AddressRepositoryMwTelemetry{
		Next: next,
	}
}

func (r *AddressRepositoryMwTelemetry) FindByIDAndContactID(ctx context.Context, db *gorm.DB, address *entity.Address, id string, contactID string) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByIDAndContactID(ctx, db, address, id, contactID)
	telemetry.RecordError(span, err)

	return err
}

func (r *AddressRepositoryMwTelemetry) FindAllByContactID(ctx context.Context, db *gorm.DB, contactID string) (entity.AddressList, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	list, err := r.Next.FindAllByContactID(ctx, db, contactID)
	telemetry.RecordError(span, err)

	return list, err
}

func (r *AddressRepositoryMwTelemetry) Create(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Create(ctx, db, entity)
	telemetry.RecordError(span, err)

	return err
}

func (r *AddressRepositoryMwTelemetry) Update(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Update(ctx, db, entity)
	telemetry.RecordError(span, err)

	return err
}

func (r *AddressRepositoryMwTelemetry) Delete(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Delete(ctx, db, entity)
	telemetry.RecordError(span, err)

	return err
}
