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

func TestAddressUsecaseImpl_Delete_Success(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	PaymentClient := &mock.PaymentClientMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		PaymentClient:     PaymentClient,
	}

	// ------------------------------------------------------- //

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectCommit()

	req := &model.DeleteAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	AddressRepository.FindByIdAndContactIdFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactId string) error {
		return nil
	}

	PaymentClient.RefundFunc = func(ctx context.Context, transactionID string) (bool, error) {
		return true, nil
	}

	AddressRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.Nil(t, err)
}

func TestAddressUsecaseImpl_Delete_Fail_FindByIdAndUserId(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	PaymentClient := &mock.PaymentClientMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		PaymentClient:     PaymentClient,
	}

	// ------------------------------------------------------- //

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectRollback()

	req := &model.DeleteAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return assert.AnError
	}

	AddressRepository.FindByIdAndContactIdFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactId string) error {
		return nil
	}

	PaymentClient.RefundFunc = func(ctx context.Context, transactionID string) (bool, error) {
		return true, nil
	}

	AddressRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_Delete_Fail_FindByIdAndContactId(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	PaymentClient := &mock.PaymentClientMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		PaymentClient:     PaymentClient,
	}

	// ------------------------------------------------------- //

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectRollback()

	req := &model.DeleteAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	AddressRepository.FindByIdAndContactIdFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactId string) error {
		return assert.AnError
	}

	PaymentClient.RefundFunc = func(ctx context.Context, transactionID string) (bool, error) {
		return true, nil
	}

	AddressRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_Delete_Fail_Refund(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	PaymentClient := &mock.PaymentClientMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		PaymentClient:     PaymentClient,
	}

	// ------------------------------------------------------- //

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectRollback()

	req := &model.DeleteAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	AddressRepository.FindByIdAndContactIdFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactId string) error {
		return nil
	}

	PaymentClient.RefundFunc = func(ctx context.Context, transactionID string) (bool, error) {
		return false, assert.AnError
	}

	AddressRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_Delete_Fail_Delete(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	PaymentClient := &mock.PaymentClientMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		PaymentClient:     PaymentClient,
	}

	// ------------------------------------------------------- //

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectRollback()

	req := &model.DeleteAddressRequest{}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	AddressRepository.FindByIdAndContactIdFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactId string) error {
		return nil
	}

	PaymentClient.RefundFunc = func(ctx context.Context, transactionID string) (bool, error) {
		return true, nil
	}

	AddressRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
