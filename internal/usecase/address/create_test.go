package address_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddressUsecaseImpl_Create_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB: gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	req := &model.CreateAddressRequest{
		UserID:    "userid1",
		ContactID: uuid.NewString(),
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
		return nil
	}

	AddressRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	expected := &model.AddressResponse{}

	// not comparing id
	res.ID = ""
	expected.ID = ""
	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestAddressUsecaseImpl_Create_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB: gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	req := &model.CreateAddressRequest{
		UserID:    "",
		ContactID: uuid.NewString(),
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
		return nil
	}

	AddressRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestAddressUsecaseImpl_Create_Fail_FindByIDAndUserID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB: gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	req := &model.CreateAddressRequest{
		UserID:    "userid1",
		ContactID: uuid.NewString(),
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
		return assert.AnError
	}

	AddressRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_Create_Fail_Create(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB: gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	req := &model.CreateAddressRequest{
		UserID:    "userid1",
		ContactID: uuid.NewString(),
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
		return nil
	}

	AddressRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return assert.AnError
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAddressUsecaseImpl_Create_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	AddressProducer := &mock.AddressProducerMock{}
	u := &address.AddressUsecaseImpl{
		DB: gormDB,
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
		AddressProducer:   AddressProducer,
	}

	// ------------------------------------------------------- //

	req := &model.CreateAddressRequest{
		UserID:    "userid1",
		ContactID: uuid.NewString(),
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
		return nil
	}

	AddressRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Address) error {
		return nil
	}

	AddressProducer.SendFunc = func(ctx context.Context, event *model.AddressEvent) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
