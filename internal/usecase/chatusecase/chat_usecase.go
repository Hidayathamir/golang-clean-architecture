package chatusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/cache"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseChat.go -pkg=mock . ChatUsecase
type ChatUsecase interface {
	Send(ctx context.Context, req dto.SendChatRequest) (dto.ChatResponse, error)
	ListConversations(ctx context.Context) (dto.ConversationsResponse, error)
	GetThread(ctx context.Context, req dto.GetThreadRequest) (dto.ThreadResponse, error)
	MarkRead(ctx context.Context, req dto.MarkChatReadRequest) error
	CountUnread(ctx context.Context) (dto.UnreadChatCountResponse, error)
}

var _ ChatUsecase = &ChatUsecaseImpl{}

type ChatUsecaseImpl struct {
	Config                *config.Config
	DB                    *gorm.DB
	ChatRepository        repository.ChatRepository
	UserRepository        repository.UserRepository
	NotificationPublisher cache.NotificationPublisher
}

func NewChatUsecase(
	Cfg *config.Config,
	DB *gorm.DB,
	chatRepository repository.ChatRepository,
	userRepository repository.UserRepository,
	notificationPublisher cache.NotificationPublisher,
) *ChatUsecaseImpl {
	return &ChatUsecaseImpl{
		Config:                Cfg,
		DB:                    DB,
		ChatRepository:        chatRepository,
		UserRepository:        userRepository,
		NotificationPublisher: notificationPublisher,
	}
}
