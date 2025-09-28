package contact_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestContactUsecaseImpl_Delete_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
	}

	// ------------------------------------------------------- //

	req := &model.DeleteContactRequest{
		UserId: "userid1",
		ID:     uuid.NewString(),
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	ContactRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.Nil(t, err)
}

func TestContactUsecaseImpl_Delete_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
	}

	// ------------------------------------------------------- //

	req := &model.DeleteContactRequest{
		UserId: "",
		ID:     uuid.NewString(),
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	ContactRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestContactUsecaseImpl_Delete_Fail_FindByIdAndUserId(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
	}

	// ------------------------------------------------------- //

	req := &model.DeleteContactRequest{
		UserId: "userid1",
		ID:     uuid.NewString(),
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return assert.AnError
	}

	ContactRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestContactUsecaseImpl_Delete_Fail_Delete(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
	}

	// ------------------------------------------------------- //

	req := &model.DeleteContactRequest{
		UserId: "userid1",
		ID:     uuid.NewString(),
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	ContactRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
