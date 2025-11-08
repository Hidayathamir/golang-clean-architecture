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

func TestUserUsecaseImpl_Current_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	u := &user.UserUsecaseImpl{
		DB: gormDB,
		 UserRepository: UserRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetUserRequest{
		ID: "userid1",
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Current(context.Background(), req)

	// ------------------------------------------------------- //

	var expected = &model.UserResponse{}

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestUserUsecaseImpl_Current_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	u := &user.UserUsecaseImpl{
		DB: gormDB,
		 UserRepository: UserRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetUserRequest{
		ID: "",
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Current(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Current_Fail_FindByID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	u := &user.UserUsecaseImpl{
		DB: gormDB,
		 UserRepository: UserRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetUserRequest{
		ID: "userid1",
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Current(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
