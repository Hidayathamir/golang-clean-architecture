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

func TestChatUsecaseImpl_ListConversations_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	chatRepository := &mock.ChatRepositoryMock{}
	userRepository := &mock.UserRepositoryMock{}
	u := &chatusecase.ChatUsecaseImpl{
		DB:             gormDB,
		ChatRepository: chatRepository,
		UserRepository: userRepository,
	}
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	chatRepository.ListConversationsFunc = func(ctx context.Context, db *gorm.DB, list *entity.ConversationList, userID int64) error {
		require.Equal(t, int64(7), userID)
		*list = entity.ConversationList{
			{PeerID: 9, LastChat: entity.Chat{ID: 1, SenderID: 9, RecipientID: 7, Content: "hi", CreatedAt: base}, LastChatAt: base},
		}
		return nil
	}
	userRepository.FindByIDsFunc = func(ctx context.Context, db *gorm.DB, userList *entity.UserList, ids []int64) error {
		require.ElementsMatch(t, []int64{9}, ids)
		*userList = entity.UserList{{ID: 9, Name: "Peer", Username: "peer"}}
		return nil
	}
	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})
	res, err := u.ListConversations(ctx)
	require.Nil(t, err)
	require.Len(t, res.Data, 1)
	require.Equal(t, "Peer", res.Data[0].PeerName)
	require.Equal(t, "peer", res.Data[0].PeerUsername)
	require.Equal(t, "hi", res.Data[0].LastChat.Content)
}
