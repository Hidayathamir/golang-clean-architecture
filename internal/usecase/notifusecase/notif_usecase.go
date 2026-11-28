package notifusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/cache"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseNotif.go -pkg=mock . NotifUsecase
type NotifUsecase interface {
	Notify(ctx context.Context, req dto.NotifyRequest) error
	List(ctx context.Context, req dto.NotificationsRequest) (dto.NotificationsResponse, error)
	MarkRead(ctx context.Context, req dto.MarkNotificationReadRequest) error
	MarkAllRead(ctx context.Context) error
	CountUnread(ctx context.Context) (dto.UnreadCountResponse, error)
}

var _ NotifUsecase = &NotifUsecaseImpl{}

type NotifUsecaseImpl struct {
	Config                 *config.Config
	DB                     *gorm.DB
	NotificationRepository repository.NotificationRepository
	UserRepository         repository.UserRepository
	NotificationPublisher  cache.NotificationPublisher
}

func NewNotifUsecase(
	Cfg *config.Config,
	DB *gorm.DB,
	notificationRepository repository.NotificationRepository,
	userRepository repository.UserRepository,
	notificationPublisher cache.NotificationPublisher,
) *NotifUsecaseImpl {
	return &NotifUsecaseImpl{
		Config:                 Cfg,
		DB:                     DB,
		NotificationRepository: notificationRepository,
		UserRepository:         userRepository,
		NotificationPublisher:  notificationPublisher,
	}
}
