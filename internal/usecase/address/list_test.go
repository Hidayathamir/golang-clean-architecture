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

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
		return nil
	}

	addresses := entity.AddressList{{}}
	AddressRepository.FindAllByContactIDFunc = func(ctx context.Context, db *gorm.DB, contactID string) (entity.AddressList, error) {
		return addresses, nil
	}

	// ------------------------------------------------------- //

	res, err := u.List(context.Background(), req)

	// ------------------------------------------------------- //

	expected := model.AddressResponseList{{}}

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestAddressUsecaseImpl_List_Fail_FindByIDAndUserID(t *testing.T) {
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

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
		return assert.AnError
	}

	var addresses entity.AddressList
	AddressRepository.FindAllByContactIDFunc = func(ctx context.Context, db *gorm.DB, contactID string) (entity.AddressList, error) {
		return addresses, nil
	}

	// ------------------------------------------------------- //

	res, err := u.List(context.Background(), req)

	// ------------------------------------------------------- //

	var expected model.AddressResponseList

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_List_Fail_FindAllByContactID(t *testing.T) {
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

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
		return nil
	}

	AddressRepository.FindAllByContactIDFunc = func(ctx context.Context, db *gorm.DB, contactID string) (entity.AddressList, error) {
		return nil, assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.List(context.Background(), req)

	// ------------------------------------------------------- //

	var expected model.AddressResponseList

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
