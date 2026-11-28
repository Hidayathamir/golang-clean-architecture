package cache

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/cachekey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/redis/go-redis/v9"
)

//go:generate moq -out=../../mock/MockCacheNotificationPublisher.go -pkg=mock . NotificationPublisher
type NotificationPublisher interface {
	Publish(ctx context.Context, userID int64, notification dto.NotificationResponse) error
}

var _ NotificationPublisher = &NotificationPublisherImpl{}

type NotificationPublisherImpl struct {
	client *redis.Client
}

func NewNotificationPublisher(client *redis.Client) NotificationPublisher {
	return &NotificationPublisherImpl{
		client: client,
	}
}

func (p *NotificationPublisherImpl) Publish(ctx context.Context, userID int64, notification dto.NotificationResponse) error {
	payload, err := json.Marshal(notification)
	if err != nil {
		return errkit.AddFuncName(err, "cache.(*NotificationPublisherImpl).Publish")
	}
	err = p.client.Publish(ctx, cachekey.NotificationChannel(userID), payload).Err()
	if err != nil {
		return errkit.AddFuncName(err, "cache.(*NotificationPublisherImpl).Publish")
	}
	return nil
}
