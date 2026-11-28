package notifusecase_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notifusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNotifUsecaseImpl_MarkRead_NotFound(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	notificationRepository := &mock.NotificationRepositoryMock{}
	u := &notifusecase.NotifUsecaseImpl{
		DB:                     gormDB,
		NotificationRepository: notificationRepository,
	}
	notificationRepository.MarkReadByIDFunc = func(ctx context.Context, db *gorm.DB, userID int64, notificationID int64) (int64, error) {
		require.Equal(t, int64(7), userID)
		require.Equal(t, int64(99), notificationID)
		return 0, nil
	}
	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})
	err := u.MarkRead(ctx, dto.MarkNotificationReadRequest{NotificationID: 99})
	require.NotNil(t, err)
}

func TestNotifUsecaseImpl_MarkRead_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	notificationRepository := &mock.NotificationRepositoryMock{}
	u := &notifusecase.NotifUsecaseImpl{
		DB:                     gormDB,
		NotificationRepository: notificationRepository,
	}
	notificationRepository.MarkReadByIDFunc = func(ctx context.Context, db *gorm.DB, userID int64, notificationID int64) (int64, error) {
		return 1, nil
	}
	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})
	err := u.MarkRead(ctx, dto.MarkNotificationReadRequest{NotificationID: 1})
	require.Nil(t, err)
}
