package chatusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/chatusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestChatUsecaseImpl_GetThread_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	chatRepository := &mock.ChatRepositoryMock{}
	u := &chatusecase.ChatUsecaseImpl{DB: gormDB, ChatRepository: chatRepository}
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	chatRepository.FindThreadFunc = func(ctx context.Context, db *gorm.DB, list *entity.ChatList, userID, peerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		require.Equal(t, int64(7), userID)
		require.Equal(t, int64(9), peerID)
		*list = entity.ChatList{
			{ID: 2, SenderID: 9, RecipientID: 7, Content: "yo", CreatedAt: base.Add(time.Hour)},
			{ID: 1, SenderID: 7, RecipientID: 9, Content: "hey", CreatedAt: base},
		}
		return nil
	}
	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})
	res, err := u.GetThread(ctx, dto.GetThreadRequest{PeerID: 9, Page: dto.PageRequest{}})
	require.Nil(t, err)
	require.Len(t, res.Data, 2)
	require.Equal(t, "yo", res.Data[0].Content)
}
