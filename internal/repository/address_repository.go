package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/AddressRepository.go -pkg=mock . AddressRepository

type AddressRepository interface {
	FindByIDAndContactID(ctx context.Context, db *gorm.DB, address *entity.Address, id string, contactID string) error
	FindAllByContactID(ctx context.Context, db *gorm.DB, contactID string) (entity.AddressList, error)
	Create(ctx context.Context, db *gorm.DB, entity *entity.Address) error
	Update(ctx context.Context, db *gorm.DB, entity *entity.Address) error
	Delete(ctx context.Context, db *gorm.DB, entity *entity.Address) error
}

var _ AddressRepository = &AddressRepositoryImpl{}

type AddressRepositoryImpl struct {
	Config *viper.Viper
}

func NewAddressRepository(cfg *viper.Viper) *AddressRepositoryImpl {
	return &AddressRepositoryImpl{
		Config: cfg,
	}
}

func (r *AddressRepositoryImpl) FindByIDAndContactID(ctx context.Context, db *gorm.DB, address *entity.Address, id string, contactID string) error {
	err := db.Where("id = ? AND contact_id = ?", id, contactID).First(address).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName("repository.(*AddressRepositoryImpl).FindByIDAndContactID", err)
	}
	return nil

}

func (r *AddressRepositoryImpl) FindAllByContactID(ctx context.Context, db *gorm.DB, contactID string) (entity.AddressList, error) {
	var addresses entity.AddressList
	if err := db.Where("contact_id = ?", contactID).Find(&addresses).Error; err != nil {
		return nil, errkit.AddFuncName("repository.(*AddressRepositoryImpl).FindAllByContactID", err)
	}
	return addresses, nil
}

func (r *AddressRepositoryImpl) Create(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	err := db.Create(entity).Error
	if err != nil {
		return errkit.AddFuncName("repository.(*AddressRepositoryImpl).Create", err)
	}
	return nil
}

func (r *AddressRepositoryImpl) Update(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	err := db.Save(entity).Error
	if err != nil {
		return errkit.AddFuncName("repository.(*AddressRepositoryImpl).Update", err)
	}
	return nil
}

func (r *AddressRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, entity *entity.Address) error {
	err := db.Delete(entity).Error
	if err != nil {
		return errkit.AddFuncName("repository.(*AddressRepositoryImpl).Delete", err)
	}
	return nil
}
