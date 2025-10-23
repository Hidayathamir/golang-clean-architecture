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

func (r *AddressRepositoryMwLogger) FindAllByContactId(ctx context.Context, db *gorm.DB, contactId string) ([]entity.Address, error) {
	addresses, err := r.Next.FindAllByContactId(ctx, db, contactId)

	fields := logrus.Fields{
		"contactId": contactId,
		"addresses": addresses,
	}
	logging.Log(ctx, fields, err)

	return addresses, err
}

func (r *AddressRepositoryMwLogger) FindByIdAndContactId(ctx context.Context, db *gorm.DB, address *entity.Address, id string, contactId string) error {
	err := r.Next.FindByIdAndContactId(ctx, db, address, id, contactId)

	fields := logrus.Fields{
		"address":   address,
		"id":        id,
		"contactId": contactId,
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
