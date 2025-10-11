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

func TestContactUsecaseImpl_Get_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetContactRequest{
		UserId: "userid1",
		ID:     uuid.NewString(),
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Get(context.Background(), req)

	// ------------------------------------------------------- //

	var expected = &model.ContactResponse{}

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestContactUsecaseImpl_Get_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetContactRequest{
		UserId: "",
		ID:     uuid.NewString(),
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Get(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.ContactResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestContactUsecaseImpl_Get_Fail_FindByIdAndUserId(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetContactRequest{
		UserId: "userid1",
		ID:     uuid.NewString(),
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Get(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.ContactResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
