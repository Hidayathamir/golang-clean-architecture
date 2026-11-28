package userusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_VerifyEmail_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	userRepository := &mock.UserRepositoryMock{}
	u := &userusecase.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: userRepository,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	var updated *entity.User
	userRepository.FindByVerificationTokenFunc = func(ctx context.Context, db *gorm.DB, user *entity.User, token string) error {
		require.Equal(t, "token-123", token)
		user.ID = 7
		user.VerificationToken = "token-123"
		return nil
	}
	userRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, user *entity.User) error {
		updated = user
		return nil
	}

	// ------------------------------------------------------- //

	err := u.VerifyEmail(context.Background(), dto.VerifyEmailRequest{Token: "token-123"})

	// ------------------------------------------------------- //

	require.Nil(t, err)
	require.NotNil(t, updated)
	require.NotNil(t, updated.EmailVerifiedAt)
}

func TestUserUsecaseImpl_VerifyEmail_Success_AlreadyVerified(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	userRepository := &mock.UserRepositoryMock{}
	u := &userusecase.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: userRepository,
		UserCache:      newUserCacheMock(t),
	}

	// ------------------------------------------------------- //

	now := time.Now()
	updateCalled := false
	userRepository.FindByVerificationTokenFunc = func(ctx context.Context, db *gorm.DB, user *entity.User, token string) error {
		user.ID = 7
		user.EmailVerifiedAt = &now
		return nil
	}
	userRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, user *entity.User) error {
		updateCalled = true
		return nil
	}

	// ------------------------------------------------------- //

	err := u.VerifyEmail(context.Background(), dto.VerifyEmailRequest{Token: "token-123"})

	// ------------------------------------------------------- //

	require.Nil(t, err)
	require.False(t, updateCalled)
}

func TestUserUsecaseImpl_VerifyEmail_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &userusecase.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: &mock.UserRepositoryMock{},
		UserCache:      newUserCacheMock(t),
	}

	err := u.VerifyEmail(context.Background(), dto.VerifyEmailRequest{Token: ""})

	require.NotNil(t, err)
}

func TestUserUsecaseImpl_VerifyEmail_Fail_NotFound(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	userRepository := &mock.UserRepositoryMock{}
	u := &userusecase.UserUsecaseImpl{
		DB:             gormDB,
		UserRepository: userRepository,
		UserCache:      newUserCacheMock(t),
	}

	userRepository.FindByVerificationTokenFunc = func(ctx context.Context, db *gorm.DB, user *entity.User, token string) error {
		return gorm.ErrRecordNotFound
	}

	err := u.VerifyEmail(context.Background(), dto.VerifyEmailRequest{Token: "missing"})

	require.NotNil(t, err)
}
