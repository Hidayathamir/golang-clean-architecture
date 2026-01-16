package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/ImageRepository.go -pkg=mock . ImageRepository

type ImageRepository interface {
	Create(ctx context.Context, db *gorm.DB, entity *entity.Image) error
}

var _ ImageRepository = &ImageRepositoryImpl{}

type ImageRepositoryImpl struct {
	Config *viper.Viper
}

func NewImageRepository(cfg *viper.Viper) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{
		Config: cfg,
	}
}

func (r *ImageRepositoryImpl) Create(ctx context.Context, db *gorm.DB, entity *entity.Image) error {
	err := db.Create(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
