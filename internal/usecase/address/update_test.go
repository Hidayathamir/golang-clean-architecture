package address_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddressUsecaseImpl_Update_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	const street = "street1"
	req := &model.UpdateAddressRequest{
		UserID:    testUserID,
		ContactID: 1,
		ID: 1,
		Street:    street,
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return nil
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return nil
	}

	AddressRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	entityAddress := &entity.Address{Street: street}
	expected := new(model.AddressResponse)
	converter.EntityAddressToModelAddressResponse(entityAddress, expected)

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestAddressUsecaseImpl_Update_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	const street = "street1"
	req := &model.UpdateAddressRequest{
		UserID:    0,
		ContactID: 1,
		ID: 1,
		Street:    street,
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return nil
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return nil
	}

	AddressRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestAddressUsecaseImpl_Update_Fail_FindByIDAndUserID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	const street = "street1"
	req := &model.UpdateAddressRequest{
		UserID:    testUserID,
		ContactID: 1,
		ID: 1,
		Street:    street,
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return assert.AnError
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return nil
	}

	AddressRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_Update_Fail_FindByIDAndContactID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	const street = "street1"
	req := &model.UpdateAddressRequest{
		UserID:    testUserID,
		ContactID: 1,
		ID: 1,
		Street:    street,
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return nil
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return assert.AnError
	}

	AddressRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_Update_Fail_Update(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	const street = "street1"
	req := &model.UpdateAddressRequest{
		UserID:    testUserID,
		ContactID: 1,
		ID: 1,
		Street:    street,
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return nil
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return nil
	}

	AddressRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return assert.AnError
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_Update_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB:                gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	const street = "street1"
	req := &model.UpdateAddressRequest{
		UserID:    testUserID,
		ContactID: 1,
		ID: 1,
		Street:    street,
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
		return nil
	}

	AddressRepository.FindByIDAndContactIDFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactID int64) error {
		return nil
	}

	AddressRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
