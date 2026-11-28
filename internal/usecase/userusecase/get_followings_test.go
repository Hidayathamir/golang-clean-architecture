package userusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_GetFollowings_Success_FirstPage(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	followRepository := &mock.FollowRepositoryMock{}
	userRepository := &mock.UserRepositoryMock{}

	u := &userusecase.UserUsecaseImpl{
		DB:               gormDB,
		FollowRepository: followRepository,
		UserRepository:   userRepository,
	}

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	followRepository.FindFollowingsByFollowerIDFunc = func(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		require.Equal(t, int64(9), followerID)
		require.Equal(t, int64(0), cursorID)
		require.Equal(t, dto.DefaultPageLimit+1, limit)
		*followList = entity.FollowList{
			{ID: 30, FollowerID: 9, FollowingID: 4, CreatedAt: base.Add(2 * time.Hour)},
			{ID: 20, FollowerID: 9, FollowingID: 5, CreatedAt: base.Add(1 * time.Hour)},
		}
		return nil
	}
	userRepository.FindByIDsFunc = func(ctx context.Context, db *gorm.DB, userList *entity.UserList, ids []int64) error {
		require.ElementsMatch(t, []int64{4, 5}, ids)
		// intentionally out of order to prove the usecase reorders.
		*userList = entity.UserList{
			{ID: 5, Username: "bob"},
			{ID: 4, Username: "alice"},
		}
		return nil
	}

	res, err := u.GetFollowings(context.Background(), dto.FollowingsRequest{UserID: 9})

	require.Nil(t, err)
	require.Len(t, res.Data, 2)
	require.Empty(t, res.NextCursor)
	require.Equal(t, int64(4), res.Data[0].ID)
	require.Equal(t, "alice", res.Data[0].Username)
	require.Equal(t, int64(5), res.Data[1].ID)
}

func TestUserUsecaseImpl_GetFollowings_Success_NextPage(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	followRepository := &mock.FollowRepositoryMock{}
	userRepository := &mock.UserRepositoryMock{}

	u := &userusecase.UserUsecaseImpl{
		DB:               gormDB,
		FollowRepository: followRepository,
		UserRepository:   userRepository,
	}

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cursor := pagination.Encode(pagination.Cursor{CreatedAt: base, ID: 100})

	followRepository.FindFollowingsByFollowerIDFunc = func(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		require.Equal(t, int64(100), cursorID)
		require.True(t, base.Equal(cursorCreatedAt))
		require.Equal(t, 3, limit)
		*followList = entity.FollowList{
			{ID: 3, FollowingID: 3, CreatedAt: base.Add(-time.Hour)},
			{ID: 2, FollowingID: 2, CreatedAt: base.Add(-2 * time.Hour)},
			{ID: 1, FollowingID: 1, CreatedAt: base.Add(-3 * time.Hour)},
		}
		return nil
	}
	userRepository.FindByIDsFunc = func(ctx context.Context, db *gorm.DB, userList *entity.UserList, ids []int64) error {
		*userList = entity.UserList{{ID: 3}, {ID: 2}}
		return nil
	}

	req := dto.FollowingsRequest{UserID: 9, Page: dto.PageRequest{Limit: 2, Cursor: cursor}}
	res, err := u.GetFollowings(context.Background(), req)

	require.Nil(t, err)
	require.Len(t, res.Data, 2)

	next, ok, err := pagination.Decode(res.NextCursor)
	require.Nil(t, err)
	require.True(t, ok)
	require.Equal(t, int64(2), next.ID)
}

func TestUserUsecaseImpl_GetFollowings_Fail_InvalidCursor(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &userusecase.UserUsecaseImpl{DB: gormDB}

	req := dto.FollowingsRequest{UserID: 9, Page: dto.PageRequest{Cursor: "not-a-valid-cursor"}}
	res, err := u.GetFollowings(context.Background(), req)

	require.Equal(t, dto.FollowingsResponse{}, res)
	require.NotNil(t, err)
}

func TestUserUsecaseImpl_GetFollowings_Fail_FollowRepository(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	followRepository := &mock.FollowRepositoryMock{}
	u := &userusecase.UserUsecaseImpl{
		DB:               gormDB,
		FollowRepository: followRepository,
		UserRepository:   &mock.UserRepositoryMock{},
	}

	followRepository.FindFollowingsByFollowerIDFunc = func(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		return assert.AnError
	}

	res, err := u.GetFollowings(context.Background(), dto.FollowingsRequest{UserID: 9})

	require.Equal(t, dto.FollowingsResponse{}, res)
	require.ErrorIs(t, err, assert.AnError)
}
