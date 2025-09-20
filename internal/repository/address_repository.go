package repository

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/pkg/errkit"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/AddressRepository.go -pkg=mock . AddressRepository

type AddressRepository interface {
	Repository[entity.Address]
	FindByIdAndContactId(db *gorm.DB, address *entity.Address, id string, contactId string) error
	FindAllByContactId(db *gorm.DB, contactId string) ([]entity.Address, error)
}

var _ AddressRepository = &AddressRepositoryImpl{}

type AddressRepositoryImpl struct {
	RepositoryImpl[entity.Address]
	Log *logrus.Logger
}

func NewAddressRepository(log *logrus.Logger) *AddressRepositoryImpl {
	return &AddressRepositoryImpl{
		RepositoryImpl: RepositoryImpl[entity.Address]{},
		Log:            log,
	}
}

func (r *AddressRepositoryImpl) FindByIdAndContactId(db *gorm.DB, address *entity.Address, id string, contactId string) error {
	err := db.Where("id = ? AND contact_id = ?", id, contactId).First(address).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil

}

func (r *AddressRepositoryImpl) FindAllByContactId(db *gorm.DB, contactId string) ([]entity.Address, error) {
	var addresses []entity.Address
	if err := db.Where("contact_id = ?", contactId).Find(&addresses).Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}
	return addresses, nil
}
