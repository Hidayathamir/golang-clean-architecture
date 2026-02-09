package notif

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseNotif.go -pkg=mock . NotifUsecase

type NotifUsecase interface {
	Notify(ctx context.Context, req *model.NotifyRequest) error
}

var _ NotifUsecase = &NotifUsecaseImpl{}

type NotifUsecaseImpl struct {
	Config *config.Config
	DB     *gorm.DB

	// repository

	// producer

	// client
}

func NewNotifUsecase(
	Cfg *config.Config,
	DB *gorm.DB,

	// repository

	// producer

	// client
) *NotifUsecaseImpl {
	return &NotifUsecaseImpl{
		Config: Cfg,
		DB:     DB,

		// repository

		// producer

		// client
	}
}
