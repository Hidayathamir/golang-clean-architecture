package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/UserRepository.go -pkg=mock . UserRepository

type UserRepository interface {
	FindByToken(ctx context.Context, db *gorm.DB, user *entity.User, token string) error
	Create(ctx context.Context, db *gorm.DB, entity *entity.User) error
	Update(ctx context.Context, db *gorm.DB, entity *entity.User) error
	CountById(ctx context.Context, db *gorm.DB, id string) (int64, error)
	FindById(ctx context.Context, db *gorm.DB, entity *entity.User, id string) error
}

var _ UserRepository = &UserRepositoryImpl{}

type UserRepositoryImpl struct {
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		Log: log,
	}
}

func (r *UserRepositoryImpl) FindByToken(ctx context.Context, db *gorm.DB, user *entity.User, token string) error {
	err := db.Where("token = ?", token).First(user).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) Create(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := db.Create(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := db.Save(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := db.Delete(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) CountById(ctx context.Context, db *gorm.DB, id string) (int64, error) {
	var total int64
	err := db.Model(new(entity.User)).Where("id = ?", id).Count(&total).Error
	if err != nil {
		return 0, errkit.AddFuncName(err)
	}
	return total, nil
}

func (r *UserRepositoryImpl) FindById(ctx context.Context, db *gorm.DB, entity *entity.User, id string) error {
	err := db.Where("id = ?", id).Take(entity).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}
