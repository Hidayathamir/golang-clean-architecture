package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockRepositoryLike.go -pkg=mock . LikeRepository

type LikeRepository interface {
	Create(ctx context.Context, db *gorm.DB, like *entity.Like) error
	FindByImageID(ctx context.Context, db *gorm.DB, likeList *entity.LikeList, imageID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error
	DeleteByUserIDAndImageID(ctx context.Context, db *gorm.DB, userID int64, imageID int64) (int64, error)
}

var _ LikeRepository = &LikeRepositoryImpl{}

type LikeRepositoryImpl struct {
	Cfg *config.Config
}

func NewLikeRepository(cfg *config.Config) *LikeRepositoryImpl {
	return &LikeRepositoryImpl{
		Cfg: cfg,
	}
}

func (r *LikeRepositoryImpl) Create(ctx context.Context, db *gorm.DB, like *entity.Like) error {
	err := db.WithContext(ctx).Create(like).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*LikeRepositoryImpl).Create")
	}
	return nil
}

func (r *LikeRepositoryImpl) FindByImageID(ctx context.Context, db *gorm.DB, likeList *entity.LikeList, imageID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	query := db.WithContext(ctx).
		Where(column.ImageID.Eq(imageID)).
		Order(column.CreatedAt.Desc()).
		Order(column.ID.Desc()).
		Limit(limit)

	if cursorID > 0 {
		whereKeyset, argCreatedAt, argID := column.KeysetBefore(cursorCreatedAt, cursorID)
		query = query.Where(whereKeyset, argCreatedAt, argID)
	}

	err := query.Find(likeList).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*LikeRepositoryImpl).FindByImageID")
	}
	return nil
}

func (r *LikeRepositoryImpl) DeleteByUserIDAndImageID(ctx context.Context, db *gorm.DB, userID int64, imageID int64) (int64, error) {
	whereUser, argUser := column.UserID.Eq(userID)
	whereImage, argImage := column.ImageID.Eq(imageID)
	result := db.WithContext(ctx).
		Where(whereUser, argUser).
		Where(whereImage, argImage).
		Delete(&entity.Like{})
	if result.Error != nil {
		return 0, errkit.AddFuncName(result.Error, "repository.(*LikeRepositoryImpl).DeleteByUserIDAndImageID")
	}
	return result.RowsAffected, nil
}
