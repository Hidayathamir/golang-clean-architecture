package userusecase_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_Unfollow_Success(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	followRepository := &mock.FollowRepositoryMock{}
	userProducer := &mock.UserProducerMock{}

	u := &userusecase.UserUsecaseImpl{
		DB:               gormDB,
		FollowRepository: followRepository,
		UserProducer:     userProducer,
	}

	followRepository.DeleteByFollowerIDAndFollowingIDFunc = func(ctx context.Context, db *gorm.DB, followerID int64, followingID int64) (int64, error) {
		require.Equal(t, int64(1), followerID)
		require.Equal(t, int64(9), followingID)
		return 1, nil
	}

	var sentEvent *dto.UserFollowedEvent
	userProducer.SendUserFollowedFunc = func(ctx context.Context, db *gorm.DB, event *dto.UserFollowedEvent) error {
		sentEvent = event
		return nil
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 1})

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectCommit()

	err := u.Unfollow(ctx, dto.UnfollowUserRequest{FollowingID: 9})

	require.Nil(t, err)
	require.NotNil(t, sentEvent)
	require.Equal(t, int64(9), sentEvent.FollowingID)
	require.Equal(t, -1, sentEvent.Count)
}

func TestUserUsecaseImpl_Unfollow_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &userusecase.UserUsecaseImpl{DB: gormDB}

	err := u.Unfollow(context.Background(), dto.UnfollowUserRequest{FollowingID: 0})

	require.NotNil(t, err)
}

func TestUserUsecaseImpl_Unfollow_Fail_NotFollowing(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	followRepository := &mock.FollowRepositoryMock{}

	u := &userusecase.UserUsecaseImpl{
		DB:               gormDB,
		FollowRepository: followRepository,
		UserProducer:     &mock.UserProducerMock{},
	}

	followRepository.DeleteByFollowerIDAndFollowingIDFunc = func(ctx context.Context, db *gorm.DB, followerID int64, followingID int64) (int64, error) {
		return 0, nil
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 1})

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectRollback()

	err := u.Unfollow(ctx, dto.UnfollowUserRequest{FollowingID: 9})

	require.NotNil(t, err)
	require.Equal(t, 404, errkit.GetHTTPError(err).HTTPCode)
}

func TestUserUsecaseImpl_Unfollow_Fail_DeleteRepository(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	followRepository := &mock.FollowRepositoryMock{}

	u := &userusecase.UserUsecaseImpl{
		DB:               gormDB,
		FollowRepository: followRepository,
		UserProducer:     &mock.UserProducerMock{},
	}

	followRepository.DeleteByFollowerIDAndFollowingIDFunc = func(ctx context.Context, db *gorm.DB, followerID int64, followingID int64) (int64, error) {
		return 0, assert.AnError
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 1})

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectRollback()

	err := u.Unfollow(ctx, dto.UnfollowUserRequest{FollowingID: 9})

	require.ErrorIs(t, err, assert.AnError)
}
