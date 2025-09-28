package address_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddressUsecaseImpl_List_Successs(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.ListAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	addresses := []entity.Address{{}}
	AddressRepository.FindAllByContactIdFunc = func(ctx context.Context, db *gorm.DB, contactId string) ([]entity.Address, error) {
		return addresses, nil
	}

	// ------------------------------------------------------- //

	res, err := u.List(context.Background(), req)

	// ------------------------------------------------------- //

	expected := []model.AddressResponse{{}}

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestAddressUsecaseImpl_List_Fail_FindByIdAndUserId(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.ListAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return assert.AnError
	}

	var addresses []entity.Address
	AddressRepository.FindAllByContactIdFunc = func(ctx context.Context, db *gorm.DB, contactId string) ([]entity.Address, error) {
		return addresses, nil
	}

	// ------------------------------------------------------- //

	res, err := u.List(context.Background(), req)

	// ------------------------------------------------------- //

	var expected []model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_List_Fail_FindAllByContactId(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.ListAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	AddressRepository.FindAllByContactIdFunc = func(ctx context.Context, db *gorm.DB, contactId string) ([]entity.Address, error) {
		return nil, assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.List(context.Background(), req)

	// ------------------------------------------------------- //

	var expected []model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
