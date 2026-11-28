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

//go:generate moq -out=../../mock/MockRepositoryFollow.go -pkg=mock . FollowRepository

type FollowRepository interface {
	Create(ctx context.Context, db *gorm.DB, follow *entity.Follow) error
	FindByFollowingID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followingID int64) error
	FindFollowersByFollowingID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followingID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error
	FindFollowingsByFollowerID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error
}

var _ FollowRepository = &FollowRepositoryImpl{}

type FollowRepositoryImpl struct {
	Cfg *config.Config
}

func NewFollowRepository(cfg *config.Config) *FollowRepositoryImpl {
	return &FollowRepositoryImpl{
		Cfg: cfg,
	}
}

func (r *FollowRepositoryImpl) Create(ctx context.Context, db *gorm.DB, follow *entity.Follow) error {
	err := db.WithContext(ctx).Create(follow).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*FollowRepositoryImpl).Create")
	}
	return nil
}

func (r *FollowRepositoryImpl) FindByFollowingID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followingID int64) error {
	err := db.WithContext(ctx).Where(column.FollowingID.Eq(followingID)).Find(followList).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*FollowRepositoryImpl).FindByFollowingID")
	}
	return nil
}

func (r *FollowRepositoryImpl) FindFollowersByFollowingID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followingID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	query := db.WithContext(ctx).
		Where(column.FollowingID.Eq(followingID)).
		Order(column.CreatedAt.Desc()).
		Order(column.ID.Desc()).
		Limit(limit)

	if cursorID > 0 {
		whereKeyset, argCreatedAt, argID := column.KeysetBefore(cursorCreatedAt, cursorID)
		query = query.Where(whereKeyset, argCreatedAt, argID)
	}

	err := query.Find(followList).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*FollowRepositoryImpl).FindFollowersByFollowingID")
	}
	return nil
}

func (r *FollowRepositoryImpl) FindFollowingsByFollowerID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	query := db.WithContext(ctx).
		Where(column.FollowerID.Eq(followerID)).
		Order(column.CreatedAt.Desc()).
		Order(column.ID.Desc()).
		Limit(limit)

	if cursorID > 0 {
		whereKeyset, argCreatedAt, argID := column.KeysetBefore(cursorCreatedAt, cursorID)
		query = query.Where(whereKeyset, argCreatedAt, argID)
	}

	err := query.Find(followList).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*FollowRepositoryImpl).FindFollowingsByFollowerID")
	}
	return nil
}
