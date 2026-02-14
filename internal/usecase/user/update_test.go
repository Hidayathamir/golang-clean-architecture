package user_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_Update_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Config:         config.NewConfig(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.UpdateUserRequest{
		ID:       1,
		Password: "password1",
		Name:     "name1",
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendUserFollowedFunc = func(ctx context.Context, event *dto.UserFollowedEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), *req)

	// ------------------------------------------------------- //

	var expected = dto.UserResponse{Name: "name1"}

	require.Equal(t, expected, res)
	require.Nil(t, err)
}

func TestUserUsecaseImpl_Update_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Config:         config.NewConfig(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.UpdateUserRequest{
		ID: 0,
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendUserFollowedFunc = func(ctx context.Context, event *dto.UserFollowedEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), *req)

	// ------------------------------------------------------- //

	var expected dto.UserResponse

	require.Equal(t, expected, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Update_Fail_FindByID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Config:         config.NewConfig(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.UpdateUserRequest{
		ID: 1,
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		return assert.AnError
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendUserFollowedFunc = func(ctx context.Context, event *dto.UserFollowedEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), *req)

	// ------------------------------------------------------- //

	var expected dto.UserResponse

	require.Equal(t, expected, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Update_Fail_Update(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Config:         config.NewConfig(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.UpdateUserRequest{
		ID: 1,
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return assert.AnError
	}

	UserProducer.SendUserFollowedFunc = func(ctx context.Context, event *dto.UserFollowedEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), *req)

	// ------------------------------------------------------- //

	var expected dto.UserResponse

	require.Equal(t, expected, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}
