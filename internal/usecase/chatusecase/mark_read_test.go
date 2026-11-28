package chatusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/chatusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestChatUsecaseImpl_MarkRead_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	chatRepository := &mock.ChatRepositoryMock{}
	u := &chatusecase.ChatUsecaseImpl{DB: gormDB, ChatRepository: chatRepository}
	chatRepository.MarkThreadReadFunc = func(ctx context.Context, db *gorm.DB, userID, peerID int64, readAt time.Time) (int64, error) {
		require.Equal(t, int64(7), userID)
		require.Equal(t, int64(9), peerID)
		return 2, nil
	}
	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})
	err := u.MarkRead(ctx, dto.MarkChatReadRequest{PeerID: 9})
	require.Nil(t, err)
}
