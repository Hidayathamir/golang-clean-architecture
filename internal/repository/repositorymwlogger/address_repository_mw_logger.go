package repositorymwlogger

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ repository.AddressRepository = &AddressRepositoryImpl{}

type AddressRepositoryImpl struct {
	logger *logrus.Logger

	next repository.AddressRepository
}

func NewAddressRepository(logger *logrus.Logger, next repository.AddressRepository) *AddressRepositoryImpl {
	return &AddressRepositoryImpl{
		logger: logger,
		next:   next,
	}
}

func (r *AddressRepositoryImpl) FindAllByContactId(ctx context.Context, db *gorm.DB, contactId string) ([]entity.Address, error) {
	addresses, err := r.next.FindAllByContactId(ctx, db, contactId)

	fields := logrus.Fields{
		"contactId": contactId,
		"addresses": addresses,
	}
	helper.Log(ctx, fields, err)

	return addresses, err
}

func (r *AddressRepositoryImpl) FindByIdAndContactId(ctx context.Context, db *gorm.DB, address *entity.Address, id string, contactId string) error {
	err := r.next.FindByIdAndContactId(ctx, db, address, id, contactId)

	fields := logrus.Fields{
		"address":   address,
		"id":        id,
		"contactId": contactId,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *AddressRepositoryImpl) Create(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	err := r.next.Create(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *AddressRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	err := r.next.Delete(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *AddressRepositoryImpl) Update(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	err := r.next.Update(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}
