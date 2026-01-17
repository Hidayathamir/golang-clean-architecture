package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/MockRepositoryImage.go -pkg=mock . ImageRepository

type ImageRepository interface {
	Create(ctx context.Context, db *gorm.DB, image *entity.Image) error
	FindByID(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error
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

func (r *ImageRepositoryImpl) Create(ctx context.Context, db *gorm.DB, image *entity.Image) error {
	err := db.WithContext(ctx).Create(image).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *ImageRepositoryImpl) FindByID(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error {
	err := db.WithContext(ctx).Where(column.ID.Eq(id)).Take(image).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}
