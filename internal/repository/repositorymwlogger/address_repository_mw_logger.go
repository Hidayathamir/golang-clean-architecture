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

	RepositoryImpl[entity.Address]
	next repository.AddressRepository
}

func NewAddressRepository(logger *logrus.Logger, next repository.AddressRepository) *AddressRepositoryImpl {
	return &AddressRepositoryImpl{
		RepositoryImpl: RepositoryImpl[entity.Address]{
			logger: logger,
			next:   next,
		},
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
