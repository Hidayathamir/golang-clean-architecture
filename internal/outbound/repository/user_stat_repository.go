package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockRepositoryUserStat.go -pkg=mock . UserStatRepository

type UserStatRepository interface {
	FindByUserID(ctx context.Context, db *gorm.DB, userStat *entity.UserStat, userID int64) error
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

func (r *UserStatRepositoryImpl) FindByUserID(ctx context.Context, db *gorm.DB, userStat *entity.UserStat, userID int64) error {
	// A user with no follow activity yet has no user_stats row; treat that as
	// zero counts rather than a not-found error.
	err := db.WithContext(ctx).Where(column.UserID.Eq(userID)).Limit(1).Find(userStat).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*UserStatRepositoryImpl).FindByUserID")
	}
	return nil
}

func (r *UserStatRepositoryImpl) IncrementFollowerCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	err := r.upsert(ctx, db, id, count, 0)
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*UserStatRepositoryImpl).IncrementFollowerCountByID")
	}
	return nil
}

func (r *UserStatRepositoryImpl) IncrementFollowingCountByID(ctx context.Context, db *gorm.DB, id int64, count int) error {
	err := r.upsert(ctx, db, id, 0, count)
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*UserStatRepositoryImpl).IncrementFollowingCountByID")
	}
	return nil
}

func (r *UserStatRepositoryImpl) IncrementFollowerCountAndFollowingCountByID(ctx context.Context, db *gorm.DB, id int64, followerCount int, followingCount int) error {
	err := r.upsert(ctx, db, id, followerCount, followingCount)
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*UserStatRepositoryImpl).IncrementFollowerCountAndFollowingCountByID")
	}
	return nil
}

// upsert applies a follower/following count delta, creating the row on first
// activity. user_stats rows are not seeded at registration, so an update-only
// increment would silently affect zero rows.
func (r *UserStatRepositoryImpl) upsert(ctx context.Context, db *gorm.DB, id int64, followerCount int, followingCount int) error {
	err := db.WithContext(ctx).
		Exec(
			`INSERT INTO `+table.UserStat+` (`+column.UserID.Str()+`, `+column.FollowerCount.Str()+`, `+column.FollowingCount.Str()+`)
			 VALUES (?, ?, ?)
			 ON CONFLICT (`+column.UserID.Str()+`) DO UPDATE SET
			   `+column.FollowerCount.Str()+` = `+table.UserStat+`.`+column.FollowerCount.Str()+` + EXCLUDED.`+column.FollowerCount.Str()+`,
			   `+column.FollowingCount.Str()+` = `+table.UserStat+`.`+column.FollowingCount.Str()+` + EXCLUDED.`+column.FollowingCount.Str()+`,
			   `+column.UpdatedAt.Str()+` = now()`,
			id, followerCount, followingCount,
		).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*UserStatRepositoryImpl).upsert")
	}
	return nil
}
