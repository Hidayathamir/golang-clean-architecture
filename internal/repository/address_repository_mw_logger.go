package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ AddressRepository = &AddressRepositoryMwLogger{}

type AddressRepositoryMwLogger struct {
	Next AddressRepository
}

func NewAddressRepositoryMwLogger(next AddressRepository) *AddressRepositoryMwLogger {
	return &AddressRepositoryMwLogger{
		Next: next,
	}
}

func (r *AddressRepositoryMwLogger) FindAllByContactID(ctx context.Context, db *gorm.DB, contactID string) (entity.AddressList, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	addresses, err := r.Next.FindAllByContactID(ctx, db, contactID)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"contactID": contactID,
		"addresses": addresses,
	}
	x.LogMw(ctx, fields, err)

	return addresses, err
}

func (r *AddressRepositoryMwLogger) FindByIDAndContactID(ctx context.Context, db *gorm.DB, address *entity.Address, id string, contactID string) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByIDAndContactID(ctx, db, address, id, contactID)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"address":   address,
		"id":        id,
		"contactID": contactID,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *AddressRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
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

func (r *AddressRepositoryMwLogger) Delete(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Delete(ctx, db, entity)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"entity": entity,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *AddressRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
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
