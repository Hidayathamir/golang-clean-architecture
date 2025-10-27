package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
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
	addresses, err := r.Next.FindAllByContactID(ctx, db, contactID)

	fields := logrus.Fields{
		"contactID": contactID,
		"addresses": addresses,
	}
	logging.Log(ctx, fields, err)

	return addresses, err
}

func (r *AddressRepositoryMwLogger) FindByIDAndContactID(ctx context.Context, db *gorm.DB, address *entity.Address, id string, contactID string) error {
	err := r.Next.FindByIDAndContactID(ctx, db, address, id, contactID)

	fields := logrus.Fields{
		"address":   address,
		"id":        id,
		"contactID": contactID,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *AddressRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	err := r.Next.Create(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *AddressRepositoryMwLogger) Delete(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	err := r.Next.Delete(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *AddressRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	err := r.Next.Update(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	logging.Log(ctx, fields, err)

	return err
}
