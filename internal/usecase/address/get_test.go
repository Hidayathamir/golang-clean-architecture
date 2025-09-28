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

func TestAddressUsecaseImpl_Get_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	AddressRepository.FindByIdAndContactIdFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactId string) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Get(context.Background(), req)

	// ------------------------------------------------------- //

	expected := &model.AddressResponse{}

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestAddressUsecaseImpl_Get_Fail_FindByIdAndUserId(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return assert.AnError
	}

	AddressRepository.FindByIdAndContactIdFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactId string) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Get(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_Get_Fail_FindByIdAndContactId(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	AddressRepository.FindByIdAndContactIdFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactId string) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Get(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
