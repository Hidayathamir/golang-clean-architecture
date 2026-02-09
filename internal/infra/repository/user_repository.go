package repository

import (
	"context"
	"errors"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockRepositoryUser.go -pkg=mock . UserRepository

type UserRepository interface {
	Create(ctx context.Context, db *gorm.DB, user *entity.User) error
	Update(ctx context.Context, db *gorm.DB, user *entity.User) error
	CountByUsername(ctx context.Context, db *gorm.DB, username string) (int64, error)
	FindByID(ctx context.Context, db *gorm.DB, user *entity.User, id int64) error
	FindByUsername(ctx context.Context, db *gorm.DB, user *entity.User, username string) error
}

var _ UserRepository = &UserRepositoryImpl{}

type UserRepositoryImpl struct {
	Cfg *config.Config
}

func NewUserRepository(cfg *config.Config) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		Cfg: cfg,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, db *gorm.DB, user *entity.User) error {
	err := db.WithContext(ctx).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = errkit.SetHTTPError(err, http.StatusConflict)
		}
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, db *gorm.DB, user *entity.User) error {
	err := db.WithContext(ctx).Save(user).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, user *entity.User) error {
	err := db.WithContext(ctx).Delete(user).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) CountByUsername(ctx context.Context, db *gorm.DB, username string) (int64, error) {
	var total int64
	err := db.WithContext(ctx).Model(new(entity.User)).Where(column.Username.Eq(username)).Count(&total).Error
	if err != nil {
		return 0, errkit.AddFuncName(err)
	}
	return total, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, db *gorm.DB, user *entity.User, id int64) error {
	err := db.WithContext(ctx).Where(column.ID.Eq(id)).Take(user).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, db *gorm.DB, user *entity.User, username string) error {
	err := db.WithContext(ctx).Where(column.Username.Eq(username)).Take(user).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}
