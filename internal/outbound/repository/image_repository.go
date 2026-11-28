package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockRepositoryImage.go -pkg=mock . ImageRepository

type ImageRepository interface {
	Create(ctx context.Context, db *gorm.DB, image *entity.Image) error
	FindByID(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error
	FindFeedByFollowerID(ctx context.Context, db *gorm.DB, imageList *entity.ImageList, followerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error
	FindByUserID(ctx context.Context, db *gorm.DB, imageList *entity.ImageList, userID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error
	IncrementCommentCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error
	IncrementLikeCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error
}

var _ ImageRepository = &ImageRepositoryImpl{}

type ImageRepositoryImpl struct {
	Cfg *config.Config
}

func NewImageRepository(cfg *config.Config) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{
		Cfg: cfg,
	}
}

func (r *ImageRepositoryImpl) Create(ctx context.Context, db *gorm.DB, image *entity.Image) error {
	err := db.WithContext(ctx).Create(image).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*ImageRepositoryImpl).Create")
	}
	return nil
}

func (r *ImageRepositoryImpl) FindByID(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error {
	err := db.WithContext(ctx).Where(column.ID.Eq(id)).Take(image).Error
	if err != nil {
		err = errkit.SetCode(err, http.StatusNotFound)
		return errkit.AddFuncName(err, "repository.(*ImageRepositoryImpl).FindByID")
	}
	return nil
}

func (r *ImageRepositoryImpl) FindFeedByFollowerID(ctx context.Context, db *gorm.DB, imageList *entity.ImageList, followerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	following := db.WithContext(ctx).
		Table(table.Follow).
		Select(column.FollowingID.Str()).
		Where(column.FollowerID.Eq(followerID))

	whereIn, argsIn := column.UserID.InSubquery(following)
	query := db.WithContext(ctx).
		Where(whereIn, argsIn).
		Order(column.CreatedAt.Desc()).
		Order(column.ID.Desc()).
		Limit(limit)

	if cursorID > 0 {
		whereKeyset, argCreatedAt, argID := column.KeysetBefore(cursorCreatedAt, cursorID)
		query = query.Where(whereKeyset, argCreatedAt, argID)
	}

	err := query.Find(imageList).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*ImageRepositoryImpl).FindFeedByFollowerID")
	}
	return nil
}

func (r *ImageRepositoryImpl) FindByUserID(ctx context.Context, db *gorm.DB, imageList *entity.ImageList, userID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	query := db.WithContext(ctx).
		Where(column.UserID.Eq(userID)).
		Order(column.CreatedAt.Desc()).
		Order(column.ID.Desc()).
		Limit(limit)

	if cursorID > 0 {
		whereKeyset, argCreatedAt, argID := column.KeysetBefore(cursorCreatedAt, cursorID)
		query = query.Where(whereKeyset, argCreatedAt, argID)
	}

	err := query.Find(imageList).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*ImageRepositoryImpl).FindByUserID")
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
		return errkit.AddFuncName(err, "repository.(*ImageRepositoryImpl).IncrementCommentCountByID")
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
		return errkit.AddFuncName(err, "repository.(*ImageRepositoryImpl).IncrementLikeCountByID")
	}
	return nil
}
