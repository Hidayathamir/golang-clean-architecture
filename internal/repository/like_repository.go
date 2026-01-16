package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/LikeRepository.go -pkg=mock . LikeRepository

type LikeRepository interface {
	Create(ctx context.Context, db *gorm.DB, like *entity.Like) error
}

var _ LikeRepository = &LikeRepositoryImpl{}

type LikeRepositoryImpl struct {
	Config *viper.Viper
}

func NewLikeRepository(cfg *viper.Viper) *LikeRepositoryImpl {
	return &LikeRepositoryImpl{
		Config: cfg,
	}
}

func (r *LikeRepositoryImpl) Create(ctx context.Context, db *gorm.DB, like *entity.Like) error {
	err := db.Create(like).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
