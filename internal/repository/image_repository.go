package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/MockRepositoryImage.go -pkg=mock . ImageRepository

type ImageRepository interface {
	Create(ctx context.Context, db *gorm.DB, image *entity.Image) error
	FindByID(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error
	IncrementCommentCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error
	IncrementLikeCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error
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

func (r *ImageRepositoryImpl) IncrementCommentCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	err := db.WithContext(ctx).
		Table(table.Image).
		Where(column.ID.Eq(id)).
		Updates(map[string]any{
			column.CommentCount.Str(): gorm.Expr(column.CommentCount.Plus(count)),
		}).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *ImageRepositoryImpl) IncrementLikeCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	err := db.WithContext(ctx).
		Table(table.Image).
		Where(column.ID.Eq(id)).
		Updates(map[string]any{
			column.LikeCount.Str(): gorm.Expr(column.LikeCount.Plus(count)),
		}).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
