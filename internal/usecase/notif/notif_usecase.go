package notif

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseNotif.go -pkg=mock . NotifUsecase

type NotifUsecase interface {
	Notify(ctx context.Context, req dto.NotifyRequest) error
}

var _ NotifUsecase = &NotifUsecaseImpl{}

type NotifUsecaseImpl struct {
	Config *config.Config
	DB     *gorm.DB

	// repository

	// producer

	// storage
}

func NewNotifUsecase(
	Cfg *config.Config,
	DB *gorm.DB,

	// repository

	// producer

	// storage
) *NotifUsecaseImpl {
	return &NotifUsecaseImpl{
		Config: Cfg,
		DB:     DB,

		// repository

		// producer

		// storage
	}
}
