package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_Create_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
	}

	// ------------------------------------------------------- //

	req := &model.RegisterUserRequest{
		Username: "user1",
		Password: "pw1",
		Name:     "name1",
	}

	UserRepository.CountByUsernameFunc = func(ctx context.Context, db *gorm.DB, username string) (int64, error) {
		return 0, nil
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendUserFollowedFunc = func(ctx context.Context, event *model.UserFollowedEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected = &model.UserResponse{
		ID:        0,
		Username:  "user1",
		Name:      "name1",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	require.Equal(t, expected, res)
	require.Nil(t, err)
}

func TestUserUsecaseImpl_Create_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
	}

	// ------------------------------------------------------- //

	req := &model.RegisterUserRequest{
		Username: "",
		Password: "pw1",
		Name:     "name1",
	}

	UserRepository.CountByUsernameFunc = func(ctx context.Context, db *gorm.DB, username string) (int64, error) {
		return 0, nil
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendUserFollowedFunc = func(ctx context.Context, event *model.UserFollowedEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	require.Equal(t, expected, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Create_Fail_Create(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
	}

	// ------------------------------------------------------- //

	req := &model.RegisterUserRequest{
		Username: "user1",
		Password: "pw1",
		Name:     "name1",
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	require.Equal(t, expected, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}