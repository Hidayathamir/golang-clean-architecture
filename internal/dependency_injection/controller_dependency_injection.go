package dependency_injection

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/websocket"
)

type Controllers struct {
	UserController         *http.UserController
	ImageController        *http.ImageController
	NotificationController *http.NotificationController
	NotificationHandler    *websocket.NotificationHandler
	NotificationSubscriber *websocket.NotificationSubscriber
}

func SetupControllers(cfg *config.Config, usecases *Usecases) *Controllers {
	userController := http.NewUserController(cfg, usecases.UserUsecase)
	imageController := http.NewImageController(cfg, usecases.ImageUsecase)
	notificationController := http.NewNotificationController(cfg, usecases.NotifUsecase)
	notificationHandler := websocket.NewNotificationHandler(usecases.NotificationHub, usecases.UserUsecase)
	notificationSubscriber := websocket.NewNotificationSubscriber(usecases.RedisClient, usecases.NotificationHub)
	return &Controllers{
		UserController:         userController,
		ImageController:        imageController,
		NotificationController: notificationController,
		NotificationHandler:    notificationHandler,
		NotificationSubscriber: notificationSubscriber,
	}
}
