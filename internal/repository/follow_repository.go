package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/MockRepositoryFollow.go -pkg=mock . FollowRepository

type FollowRepository interface {
	Create(ctx context.Context, db *gorm.DB, follow *entity.Follow) error
	FindByFollowingID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followingID int64) error
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

func (r *FollowRepositoryImpl) Create(ctx context.Context, db *gorm.DB, follow *entity.Follow) error {
	err := db.WithContext(ctx).Create(follow).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *FollowRepositoryImpl) FindByFollowingID(ctx context.Context, db *gorm.DB, followList *entity.FollowList, followingID int64) error {
	err := db.WithContext(ctx).Where(column.FollowingID.Eq(followingID)).Find(followList).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
