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

func TestContactUsecaseImpl_Search_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.SearchContactRequest{
		UserId: "userid1",
		Page:   1,
		Size:   1,
	}

	ContactRepository.SearchFunc = func(ctx context.Context, db *gorm.DB, req *model.SearchContactRequest) ([]entity.Contact, int64, error) {
		return []entity.Contact{{}}, 12, nil
	}

	// ------------------------------------------------------- //

	res, total, err := u.Search(context.Background(), req)

	// ------------------------------------------------------- //

	var expected = []model.ContactResponse{{}}

	assert.Equal(t, expected, res)
	assert.Equal(t, int64(12), total)
	assert.Nil(t, err)
}

func TestContactUsecaseImpl_Search_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.SearchContactRequest{
		UserId: "",
		Page:   1,
		Size:   1,
	}

	ContactRepository.SearchFunc = func(ctx context.Context, db *gorm.DB, req *model.SearchContactRequest) ([]entity.Contact, int64, error) {
		return []entity.Contact{{}}, 12, nil
	}

	// ------------------------------------------------------- //

	res, total, err := u.Search(context.Background(), req)

	// ------------------------------------------------------- //

	var expected []model.ContactResponse

	assert.Equal(t, expected, res)
	assert.Equal(t, int64(0), total)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestContactUsecaseImpl_Search_Fail_Search(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	req := &model.SearchContactRequest{
		UserId: "userid1",
		Page:   1,
		Size:   1,
	}

	ContactRepository.SearchFunc = func(ctx context.Context, db *gorm.DB, req *model.SearchContactRequest) ([]entity.Contact, int64, error) {
		return nil, 0, assert.AnError
	}

	// ------------------------------------------------------- //

	res, total, err := u.Search(context.Background(), req)

	// ------------------------------------------------------- //

	var expected []model.ContactResponse

	assert.Equal(t, expected, res)
	assert.Equal(t, int64(0), total)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
