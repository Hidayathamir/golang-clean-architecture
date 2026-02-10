package user_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_Current_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: UserRepository,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.GetUserRequest{
		ID: 1,
	}

	u.UserCache = &mock.UserCacheMock{
		GetFunc: func(ctx context.Context, id int64) (*entity.User, error) {
			return nil, nil
		},
		SetFunc: func(ctx context.Context, user *entity.User) error {
			return nil
		},
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Current(context.Background(), req)

	// ------------------------------------------------------- //

	var expected = &dto.UserResponse{}

	require.Equal(t, expected, res)
	require.Nil(t, err)
}

func TestUserUsecaseImpl_Current_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: UserRepository,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.GetUserRequest{
		ID: 0,
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Current(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *dto.UserResponse

	require.Equal(t, expected, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Current_Fail_FindByID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: UserRepository,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.GetUserRequest{
		ID: 1,
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Current(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *dto.UserResponse

	require.Equal(t, expected, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}
