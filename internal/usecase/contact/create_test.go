package contact_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestContactUsecaseImpl_Create_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
	}

	// ------------------------------------------------------- //

	req := &model.CreateContactRequest{
		UserID:    testUserID,
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	ContactRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	ContactProducer.SendFunc = func(ctx context.Context, event *model.ContactEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	expected := &model.ContactResponse{
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	// not comparing id
	res.ID = 0
	expected.ID = 0
	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestContactUsecaseImpl_Create_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
	}

	// ------------------------------------------------------- //

	req := &model.CreateContactRequest{
		UserID:    0,
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	ContactRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	ContactProducer.SendFunc = func(ctx context.Context, event *model.ContactEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.ContactResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestContactUsecaseImpl_Create_Fail_Create(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
	}

	// ------------------------------------------------------- //

	req := &model.CreateContactRequest{
		UserID:    testUserID,
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	ContactRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return assert.AnError
	}

	ContactProducer.SendFunc = func(ctx context.Context, event *model.ContactEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.ContactResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestContactUsecaseImpl_Create_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
	}

	// ------------------------------------------------------- //

	req := &model.CreateContactRequest{
		UserID:    testUserID,
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	ContactRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	ContactProducer.SendFunc = func(ctx context.Context, event *model.ContactEvent) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.ContactResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
