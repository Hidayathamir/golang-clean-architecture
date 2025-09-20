package repositorymwlogger

import (
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

func (r *AddressRepositoryImpl) FindAllByContactId(db *gorm.DB, contactId string) ([]entity.Address, error) {
	addresses, err := r.next.FindAllByContactId(db, contactId)

	fields := logrus.Fields{
		"contactId": contactId,
		"addresses": addresses,
	}
	helper.Log(r.logger, fields, err)

	return addresses, err
}

func (r *AddressRepositoryImpl) FindByIdAndContactId(db *gorm.DB, address *entity.Address, id string, contactId string) error {
	err := r.next.FindByIdAndContactId(db, address, id, contactId)

	fields := logrus.Fields{
		"address":   address,
		"id":        id,
		"contactId": contactId,
	}
	helper.Log(r.logger, fields, err)

	return err
}

func (r *AddressRepositoryImpl) Create(db *gorm.DB, entity *entity.Address) error {
	err := r.next.Create(db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(r.logger, fields, err)

	return err
}

func (r *AddressRepositoryImpl) Delete(db *gorm.DB, entity *entity.Address) error {
	err := r.next.Delete(db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(r.logger, fields, err)

	return err
}

func (r *AddressRepositoryImpl) Update(db *gorm.DB, entity *entity.Address) error {
	err := r.next.Update(db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(r.logger, fields, err)

	return err
}
