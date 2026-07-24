package dependency_injection

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/idempotencyusecase"
)

type Consumers struct {
	ImageConsumer      *messaging.ImageConsumer
	NotifConsumer      *messaging.NotifConsumer
	UserConsumer       *messaging.UserConsumer
	IdempotencyUsecase idempotencyusecase.IdempotencyUsecase
}

func SetupConsumers(cfg *config.Config, usecases *Usecases) *Consumers {
	imageConsumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	notifConsumer := messaging.NewNotifConsumer(usecases.NotifUsecase)
	userConsumer := messaging.NewUserConsumer(usecases.UserUsecase)

	return &Consumers{
		ImageConsumer:      imageConsumer,
		NotifConsumer:      notifConsumer,
		UserConsumer:       userConsumer,
		IdempotencyUsecase: usecases.IdempotencyUsecase,
	}
}
