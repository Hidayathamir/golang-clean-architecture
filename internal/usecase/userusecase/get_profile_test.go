package userusecase_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_GetProfile_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	userRepository := &mock.UserRepositoryMock{}
	userStatRepository := &mock.UserStatRepositoryMock{}

	u := &userusecase.UserUsecaseImpl{
		DB:                 gormDB,
		UserRepository:     userRepository,
		UserStatRepository: userStatRepository,
	}

	userRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, user *entity.User, id int64) error {
		require.Equal(t, int64(3), id)
		*user = entity.User{ID: 3, Username: "alice", Name: "Alice"}
		return nil
	}
	userStatRepository.FindByUserIDFunc = func(ctx context.Context, db *gorm.DB, userStat *entity.UserStat, userID int64) error {
		*userStat = entity.UserStat{FollowerCount: 10, FollowingCount: 4}
		return nil
	}

	res, err := u.GetProfile(context.Background(), dto.GetProfileRequest{ID: 3})

	require.Nil(t, err)
	require.Equal(t, int64(3), res.ID)
	require.Equal(t, "alice", res.Username)
	require.Equal(t, 10, res.FollowerCount)
	require.Equal(t, 4, res.FollowingCount)
}

func TestUserUsecaseImpl_GetProfile_Success_NoStatRow(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	userRepository := &mock.UserRepositoryMock{}
	userStatRepository := &mock.UserStatRepositoryMock{}

	u := &userusecase.UserUsecaseImpl{
		DB:                 gormDB,
		UserRepository:     userRepository,
		UserStatRepository: userStatRepository,
	}

	userRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, user *entity.User, id int64) error {
		*user = entity.User{ID: 3, Username: "alice"}
		return nil
	}
	// userStat is left zero-valued: a user with no user_stats row.
	userStatRepository.FindByUserIDFunc = func(ctx context.Context, db *gorm.DB, userStat *entity.UserStat, userID int64) error {
		return nil
	}

	res, err := u.GetProfile(context.Background(), dto.GetProfileRequest{ID: 3})

	require.Nil(t, err)
	require.Equal(t, 0, res.FollowerCount)
	require.Equal(t, 0, res.FollowingCount)
}

func TestUserUsecaseImpl_GetProfile_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &userusecase.UserUsecaseImpl{
		DB:                 gormDB,
		UserRepository:     &mock.UserRepositoryMock{},
		UserStatRepository: &mock.UserStatRepositoryMock{},
	}

	res, err := u.GetProfile(context.Background(), dto.GetProfileRequest{ID: 0})

	require.Equal(t, dto.ProfileResponse{}, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_GetProfile_Fail_FindByID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	userRepository := &mock.UserRepositoryMock{}
	u := &userusecase.UserUsecaseImpl{
		DB:                 gormDB,
		UserRepository:     userRepository,
		UserStatRepository: &mock.UserStatRepositoryMock{},
	}

	userRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, user *entity.User, id int64) error {
		return assert.AnError
	}

	res, err := u.GetProfile(context.Background(), dto.GetProfileRequest{ID: 1})

	require.Equal(t, dto.ProfileResponse{}, res)
	require.ErrorIs(t, err, assert.AnError)
}
