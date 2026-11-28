package chatusecase_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/chatusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestChatUsecaseImpl_CountUnread_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	chatRepository := &mock.ChatRepositoryMock{}
	u := &chatusecase.ChatUsecaseImpl{DB: gormDB, ChatRepository: chatRepository}
	chatRepository.CountUnreadFunc = func(ctx context.Context, db *gorm.DB, userID int64) (int64, error) {
		require.Equal(t, int64(7), userID)
		return 3, nil
	}
	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})
	res, err := u.CountUnread(ctx)
	require.Nil(t, err)
	require.Equal(t, int64(3), res.UnreadCount)
}
