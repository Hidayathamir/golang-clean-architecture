package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/cachekey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/redis/go-redis/v9"
)

// NotificationSubscriber fans Redis-published notifications out to the local
// hub. This is what lets a notification produced in the worker process reach a
// WebSocket held by the webserver process.
type NotificationSubscriber struct {
	client *redis.Client
	hub    *Hub
}

func NewNotificationSubscriber(client *redis.Client, hub *Hub) *NotificationSubscriber {
	return &NotificationSubscriber{
		client: client,
		hub:    hub,
	}
}

// Run blocks until ctx is cancelled, forwarding every message on the
// notification channel pattern to the hub connection for that user.
func (s *NotificationSubscriber) Run(ctx context.Context) {
	pubsub := s.client.PSubscribe(ctx, cachekey.NotificationChannelPattern)
	defer func() {
		err := pubsub.Close()
		logkit.LogIfErr(err)
	}()
	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			s.handle(ctx, msg)
		}
	}
}

func (s *NotificationSubscriber) handle(ctx context.Context, msg *redis.Message) {
	userID, err := userIDFromChannel(msg.Channel)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error("invalid notification channel")
		return
	}
	notification := dto.NotificationResponse{}
	err = json.Unmarshal([]byte(msg.Payload), &notification)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error("invalid notification payload")
		return
	}
	s.hub.Send(userID, []byte(msg.Payload))
}

func userIDFromChannel(channel string) (int64, error) {
	idx := strings.LastIndex(channel, ":")
	if idx < 0 {
		return 0, errkit.AddFuncName(errors.New("invalid notification channel"), "websocket.userIDFromChannel")
	}
	userID, err := strconv.ParseInt(channel[idx+1:], 10, 64)
	if err != nil {
		return 0, errkit.AddFuncName(err, "websocket.userIDFromChannel")
	}
	return userID, nil
}
