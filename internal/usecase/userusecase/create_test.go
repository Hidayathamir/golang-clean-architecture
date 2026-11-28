package userusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_Create_Success(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	EmailProducer := &mock.EmailProducerMock{}
	u := &userusecase.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: UserRepository,
		EmailProducer:  EmailProducer,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.RegisterUserRequest{
		Username: "user1",
		Password: "pw1",
		Name:     "name1",
		Email:    "user1@example.com",
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		entityMoqParam.ID = 7
		return nil
	}

	var sentEvent *dto.EmailVerificationEvent
	EmailProducer.SendVerificationEmailFunc = func(ctx context.Context, db *gorm.DB, event *dto.EmailVerificationEvent) error {
		sentEvent = event
		return nil
	}

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectCommit()

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), *req)

	// ------------------------------------------------------- //

	var expected = dto.UserResponse{
		ID:        7,
		Username:  "user1",
		Name:      "name1",
		Email:     "user1@example.com",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	require.Nil(t, err)
	require.Equal(t, expected, res)

	require.NotNil(t, sentEvent)
	require.Equal(t, int64(7), sentEvent.UserID)
	require.Equal(t, "user1@example.com", sentEvent.Email)
	require.NotEmpty(t, sentEvent.Token)
}

func TestUserUsecaseImpl_Create_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	EmailProducer := &mock.EmailProducerMock{}
	u := &userusecase.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: UserRepository,
		EmailProducer:  EmailProducer,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.RegisterUserRequest{
		Username: "",
		Password: "pw1",
		Name:     "name1",
		Email:    "not-an-email",
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), *req)

	// ------------------------------------------------------- //

	var expected dto.UserResponse

	require.Equal(t, expected, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Create_Fail_Create(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	EmailProducer := &mock.EmailProducerMock{}
	u := &userusecase.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: UserRepository,
		EmailProducer:  EmailProducer,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.RegisterUserRequest{
		Username: "user1",
		Password: "pw1",
		Name:     "name1",
		Email:    "user1@example.com",
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return assert.AnError
	}

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectRollback()

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), *req)

	// ------------------------------------------------------- //

	var expected dto.UserResponse

	require.Equal(t, expected, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Create_Fail_SendVerificationEmail(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	EmailProducer := &mock.EmailProducerMock{}
	u := &userusecase.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: UserRepository,
		EmailProducer:  EmailProducer,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	req := &dto.RegisterUserRequest{
		Username: "user1",
		Password: "pw1",
		Name:     "name1",
		Email:    "user1@example.com",
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	EmailProducer.SendVerificationEmailFunc = func(ctx context.Context, db *gorm.DB, event *dto.EmailVerificationEvent) error {
		return assert.AnError
	}

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectRollback()

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), *req)

	// ------------------------------------------------------- //

	var expected dto.UserResponse

	require.Equal(t, expected, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}
