package repository

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/pkg/errkit"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/AddressRepository.go -pkg=mock . AddressRepository

type AddressRepository interface {
	FindByIdAndContactId(db *gorm.DB, address *entity.Address, id string, contactId string) error
	FindAllByContactId(db *gorm.DB, contactId string) ([]entity.Address, error)
	Create(db *gorm.DB, entity *entity.Address) error
	Update(db *gorm.DB, entity *entity.Address) error
	Delete(db *gorm.DB, entity *entity.Address) error
}

var _ AddressRepository = &AddressRepositoryImpl{}

type AddressRepositoryImpl struct {
	Log *logrus.Logger
}

func NewAddressRepository(log *logrus.Logger) *AddressRepositoryImpl {
	return &AddressRepositoryImpl{
		Log: log,
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

func (r *AddressRepositoryImpl) Create(db *gorm.DB, entity *entity.Address) error {
	err := db.Create(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *AddressRepositoryImpl) Update(db *gorm.DB, entity *entity.Address) error {
	err := db.Save(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *AddressRepositoryImpl) Delete(db *gorm.DB, entity *entity.Address) error {
	err := db.Delete(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
