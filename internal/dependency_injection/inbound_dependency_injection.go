package dependency_injection

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/websocket"
	"github.com/redis/go-redis/v9"
)

type Controllers struct {
	UserController         *http.UserController
	ImageController        *http.ImageController
	NotificationController *http.NotificationController
	NotificationHandler    *websocket.NotificationHandler
	NotificationSubscriber *websocket.NotificationSubscriber
	ChatController         *http.ChatController
}

type Hubs struct {
	NotificationHub *websocket.Hub
}

func SetupHubs() *Hubs {
	return &Hubs{
		NotificationHub: websocket.NewHub(),
	}
}

func SetupControllers(cfg *config.Config, redisClient *redis.Client, hubs *Hubs, usecases *Usecases) *Controllers {
	userController := http.NewUserController(cfg, usecases.UserUsecase)
	imageController := http.NewImageController(cfg, usecases.ImageUsecase)
	notificationController := http.NewNotificationController(cfg, usecases.NotifUsecase)
	notificationHandler := websocket.NewNotificationHandler(hubs.NotificationHub, usecases.UserUsecase)
	notificationSubscriber := websocket.NewNotificationSubscriber(redisClient, hubs.NotificationHub)
	chatController := http.NewChatController(cfg, usecases.ChatUsecase)
	return &Controllers{
		UserController:         userController,
		ImageController:        imageController,
		NotificationController: notificationController,
		NotificationHandler:    notificationHandler,
		NotificationSubscriber: notificationSubscriber,
		ChatController:         chatController,
	}
}

type Consumers struct {
	ImageConsumer *messaging.ImageConsumer
	NotifConsumer *messaging.NotifConsumer
	UserConsumer  *messaging.UserConsumer
}

func SetupConsumers(cfg *config.Config, usecases *Usecases) *Consumers {
	imageConsumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	notifConsumer := messaging.NewNotifConsumer(usecases.NotifUsecase)
	userConsumer := messaging.NewUserConsumer(usecases.UserUsecase)

	return &Consumers{
		ImageConsumer: imageConsumer,
		NotifConsumer: notifConsumer,
		UserConsumer:  userConsumer,
	}
}
