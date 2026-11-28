package userusecase_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"github.com/stretchr/testify/require"
)

func TestUserUsecaseImpl_NotifyUserBeingFollowed_EmitsNotif(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	notifUsecase := &mock.NotifUsecaseMock{}
	u := &userusecase.UserUsecaseImpl{
		DB:           gormDB,
		NotifUsecase: notifUsecase,
	}
	var got dto.NotifyRequest
	notifUsecase.NotifyFunc = func(ctx context.Context, req dto.NotifyRequest) error {
		got = req
		return nil
	}
	err := u.NotifyUserBeingFollowed(context.Background(), dto.NotifyUserBeingFollowedRequest{
		FollowerID:  4,
		FollowingID: 9,
	})
	require.Nil(t, err)
	require.Equal(t, int64(9), got.UserID)
	require.Equal(t, int64(4), got.ActorID)
	require.Equal(t, entity.NotificationTypeFollow, got.Type)
}
