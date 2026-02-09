package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockRepositoryUserStat.go -pkg=mock . UserStatRepository

type UserStatRepository interface {
	IncrementFollowerCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error
	IncrementFollowingCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error
	IncrementFollowerCountAndFollowingCountByID(ctx context.Context, db *gorm.DB, id int64, followerCount int, followingCount int) error
}

var _ UserStatRepository = &UserStatRepositoryImpl{}

type UserStatRepositoryImpl struct {
	Cfg *config.Config
}

func NewUserStatRepository(cfg *config.Config) *UserStatRepositoryImpl {
	return &UserStatRepositoryImpl{
		Cfg: cfg,
	}
}

func (r *UserStatRepositoryImpl) IncrementFollowerCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	err := db.WithContext(ctx).
		Table(table.UserStat).
		Where(column.UserID.Eq(id)).
		Updates(map[string]any{
			column.FollowerCount.Str(): gorm.Expr(column.FollowerCount.Plus(count)),
		}).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserStatRepositoryImpl) IncrementFollowingCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	err := db.WithContext(ctx).
		Table(table.UserStat).
		Where(column.UserID.Eq(id)).
		Updates(map[string]any{
			column.FollowingCount.Str(): gorm.Expr(column.FollowingCount.Plus(count)),
		}).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserStatRepositoryImpl) IncrementFollowerCountAndFollowingCountByID(ctx context.Context, db *gorm.DB, id int64, followerCount int, followingCount int) error {
	err := db.WithContext(ctx).
		Table(table.UserStat).
		Where(column.UserID.Eq(id)).
		Updates(map[string]any{
			column.FollowerCount.Str():  gorm.Expr(column.FollowerCount.Plus(followerCount)),
			column.FollowingCount.Str(): gorm.Expr(column.FollowingCount.Plus(followingCount)),
		}).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
