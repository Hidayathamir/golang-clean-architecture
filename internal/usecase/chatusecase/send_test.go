package chatusecase_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/chatusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestChatUsecaseImpl_Send_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	chatRepository := &mock.ChatRepositoryMock{}
	userRepository := &mock.UserRepositoryMock{}
	publisher := &mock.NotificationPublisherMock{}
	u := &chatusecase.ChatUsecaseImpl{
		DB:                    gormDB,
		ChatRepository:        chatRepository,
		UserRepository:        userRepository,
		NotificationPublisher: publisher,
	}
	userRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, user *entity.User, id int64) error {
		user.ID = id
		user.Name = "Peer"
		return nil
	}
	var created *entity.Chat
	chatRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, chat *entity.Chat) error {
		chat.ID = 1
		created = chat
		return nil
	}
	var published dto.NotificationResponse
	publisher.PublishFunc = func(ctx context.Context, userID int64, notification dto.NotificationResponse) error {
		require.Equal(t, int64(9), userID)
		published = notification
		return nil
	}
	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})
	res, err := u.Send(ctx, dto.SendChatRequest{RecipientID: 9, Content: "hello"})
	require.Nil(t, err)
	require.Equal(t, int64(7), created.SenderID)
	require.Equal(t, int64(9), created.RecipientID)
	require.Equal(t, "hello", res.Content)
	require.Equal(t, int64(1), published.ID)
}

func TestChatUsecaseImpl_Send_SelfFails(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &chatusecase.ChatUsecaseImpl{
		DB:                    gormDB,
		ChatRepository:        &mock.ChatRepositoryMock{},
		UserRepository:        &mock.UserRepositoryMock{},
		NotificationPublisher: &mock.NotificationPublisherMock{},
	}
	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})
	_, err := u.Send(ctx, dto.SendChatRequest{RecipientID: 7, Content: "hi"})
	require.NotNil(t, err)
}
