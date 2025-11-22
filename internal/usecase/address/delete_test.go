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

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return nil
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return nil
	}

	PaymentClient.RefundFunc = func(ctx context.Context, req model.PaymentRefundRequest) (model.PaymentRefundResponse, error) {
		return model.PaymentRefundResponse{Success: true}, nil
	}

	AddressRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.Nil(t, err)
}

func TestAddressUsecaseImpl_Delete_Fail_FindByIDAndUserID(t *testing.T) {
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

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return assert.AnError
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return nil
	}

	PaymentClient.RefundFunc = func(ctx context.Context, req model.PaymentRefundRequest) (model.PaymentRefundResponse, error) {
		return model.PaymentRefundResponse{Success: true}, nil
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

func TestAddressUsecaseImpl_Delete_Fail_FindByIDAndContactID(t *testing.T) {
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

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return nil
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return assert.AnError
	}

	PaymentClient.RefundFunc = func(ctx context.Context, req model.PaymentRefundRequest) (model.PaymentRefundResponse, error) {
		return model.PaymentRefundResponse{Success: true}, nil
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

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return nil
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return nil
	}

	PaymentClient.RefundFunc = func(ctx context.Context, req model.PaymentRefundRequest) (model.PaymentRefundResponse, error) {
		return model.PaymentRefundResponse{Success: false}, assert.AnError
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

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return nil
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return nil
	}

	PaymentClient.RefundFunc = func(ctx context.Context, req model.PaymentRefundRequest) (model.PaymentRefundResponse, error) {
		return model.PaymentRefundResponse{Success: true}, nil
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
