package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/User2Repository.go -pkg=mock . User2Repository

type User2Repository interface {
	Create(ctx context.Context, db *gorm.DB, user *entity.User2) error
	FindByEmail(ctx context.Context, db *gorm.DB, user *entity.User2, email string) error
	FindByID(ctx context.Context, db *gorm.DB, user *entity.User2, id string) error
	Update(ctx context.Context, db *gorm.DB, user *entity.User2) error
}

var _ User2Repository = &User2RepositoryImpl{}

type User2RepositoryImpl struct {
	Config *viper.Viper
	Log    *logrus.Logger
}

func NewUser2Repository(cfg *viper.Viper, log *logrus.Logger) *User2RepositoryImpl {
	return &User2RepositoryImpl{
		Config: cfg,
		Log:    log,
	}
}

func (r *User2RepositoryImpl) Create(ctx context.Context, db *gorm.DB, user *entity.User2) error {
	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		return errkit.AddFuncName("repository.(*User2RepositoryImpl).Create", err)
	}
	return nil
}

func (r *User2RepositoryImpl) FindByEmail(ctx context.Context, db *gorm.DB, user *entity.User2, email string) error {
	err := db.WithContext(ctx).Where("email = ?", email).Take(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errkit.NotFound(err)
		}
		return errkit.AddFuncName("repository.(*User2RepositoryImpl).FindByEmail", err)
	}
	return nil
}

func (r *User2RepositoryImpl) FindByID(ctx context.Context, db *gorm.DB, user *entity.User2, id string) error {
	err := db.WithContext(ctx).Where("id = ?", id).Take(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errkit.NotFound(err)
		}
		return errkit.AddFuncName("repository.(*User2RepositoryImpl).FindByID", err)
	}
	return nil
}

func (r *User2RepositoryImpl) Update(ctx context.Context, db *gorm.DB, user *entity.User2) error {
	if err := db.WithContext(ctx).Save(user).Error; err != nil {
		return errkit.AddFuncName("repository.(*User2RepositoryImpl).Update", err)
	}
	return nil
}
