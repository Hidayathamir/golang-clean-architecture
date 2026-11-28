package notifusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notifusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNotifUsecaseImpl_List_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	notificationRepository := &mock.NotificationRepositoryMock{}
	userRepository := &mock.UserRepositoryMock{}
	u := &notifusecase.NotifUsecaseImpl{
		DB:                     gormDB,
		NotificationRepository: notificationRepository,
		UserRepository:         userRepository,
	}
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	notificationRepository.FindByUserIDFunc = func(ctx context.Context, db *gorm.DB, list *entity.NotificationList, userID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		require.Equal(t, int64(7), userID)
		*list = entity.NotificationList{
			{ID: 2, UserID: 7, ActorID: 4, Type: entity.NotificationTypeLike, ImageID: 10, CreatedAt: base.Add(time.Hour)},
			{ID: 1, UserID: 7, ActorID: 5, Type: entity.NotificationTypeFollow, CreatedAt: base},
		}
		return nil
	}
	userRepository.FindByIDsFunc = func(ctx context.Context, db *gorm.DB, userList *entity.UserList, ids []int64) error {
		require.ElementsMatch(t, []int64{4, 5}, ids)
		*userList = entity.UserList{
			{ID: 5, Name: "Bob"},
			{ID: 4, Name: "Alice"},
		}
		return nil
	}
	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})
	res, err := u.List(ctx, dto.NotificationsRequest{})
	require.Nil(t, err)
	require.Len(t, res.Data, 2)
	require.Equal(t, "Alice", res.Data[0].ActorName)
	require.Equal(t, entity.NotificationTypeLike, res.Data[0].Type)
	require.Equal(t, "Bob", res.Data[1].ActorName)
}
