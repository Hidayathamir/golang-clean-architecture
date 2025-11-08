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
		DB: gormDB,
		 ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetContactRequest{
		UserID: "userid1",
		ID:     uuid.NewString(),
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
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
		DB: gormDB,
		 ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetContactRequest{
		UserID: "",
		ID:     uuid.NewString(),
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
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

func TestContactUsecaseImpl_Get_Fail_FindByIDAndUserID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &contact.ContactUsecaseImpl{
		DB: gormDB,
		 ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetContactRequest{
		UserID: "userid1",
		ID:     uuid.NewString(),
	}

	ContactRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userID string) error {
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
