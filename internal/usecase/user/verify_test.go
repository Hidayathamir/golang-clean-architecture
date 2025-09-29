package user_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_Verify_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
	}

	// ------------------------------------------------------- //

	req := &model.VerifyUserRequest{
		Token: "token1",
	}

	UserRepository.FindByTokenFunc = func(ctx context.Context, db *gorm.DB, user *entity.User, token string) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Verify(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.Auth = &model.Auth{}

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestUserUsecaseImpl_Verify_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
	}

	// ------------------------------------------------------- //

	req := &model.VerifyUserRequest{
		Token: "",
	}

	UserRepository.FindByTokenFunc = func(ctx context.Context, db *gorm.DB, user *entity.User, token string) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Verify(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.Auth

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Verify_FindByToken(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
	}

	// ------------------------------------------------------- //

	req := &model.VerifyUserRequest{
		Token: "token1",
	}

	UserRepository.FindByTokenFunc = func(ctx context.Context, db *gorm.DB, user *entity.User, token string) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Verify(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.Auth

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
