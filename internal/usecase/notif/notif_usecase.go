package notif

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseNotif.go -pkg=mock . NotifUsecase

type NotifUsecase interface {
	Notify(ctx context.Context, req *model.NotifyRequest) error
}

var _ NotifUsecase = &NotifUsecaseImpl{}

type NotifUsecaseImpl struct {
	Config *viper.Viper
	DB     *gorm.DB

	// repository

	// producer

	// client
}

func NewNotifUsecase(
	Config *viper.Viper,
	DB *gorm.DB,

	// repository

	// producer

	// client
) *NotifUsecaseImpl {
	return &NotifUsecaseImpl{
		Config: Config,
		DB:     DB,

		// repository

		// producer

		// client
	}
}
