package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/FollowRepository.go -pkg=mock . FollowRepository

type FollowRepository interface {
	Create(ctx context.Context, db *gorm.DB, entity *entity.Follow) error
}

var _ FollowRepository = &FollowRepositoryImpl{}

type FollowRepositoryImpl struct {
	Config *viper.Viper
}

func NewFollowRepository(cfg *viper.Viper) *FollowRepositoryImpl {
	return &FollowRepositoryImpl{
		Config: cfg,
	}
}

func (r *FollowRepositoryImpl) Create(ctx context.Context, db *gorm.DB, entity *entity.Follow) error {
	err := db.Create(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
