package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/UserRepository.go -pkg=mock . UserRepository

type UserRepository interface {
	Create(ctx context.Context, db *gorm.DB, entity *entity.User) error
	Update(ctx context.Context, db *gorm.DB, entity *entity.User) error
	CountByUsername(ctx context.Context, db *gorm.DB, username string) (int64, error)
	FindByID(ctx context.Context, db *gorm.DB, entity *entity.User, id int64) error
	FindByUsername(ctx context.Context, db *gorm.DB, entity *entity.User, username string) error
}

var _ UserRepository = &UserRepositoryImpl{}

type UserRepositoryImpl struct {
	Config *viper.Viper
}

func NewUserRepository(cfg *viper.Viper) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		Config: cfg,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := db.Create(entity).Error
	if err != nil {
		return errkit.AddFuncName("repository.(*UserRepositoryImpl).Create", err)
	}
	return nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := db.Save(entity).Error
	if err != nil {
		return errkit.AddFuncName("repository.(*UserRepositoryImpl).Update", err)
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := db.Delete(entity).Error
	if err != nil {
		return errkit.AddFuncName("repository.(*UserRepositoryImpl).Delete", err)
	}
	return nil
}

func (r *UserRepositoryImpl) CountByUsername(ctx context.Context, db *gorm.DB, username string) (int64, error) {
	var total int64
	err := db.Model(new(entity.User)).Where("username = ?", username).Count(&total).Error
	if err != nil {
		return 0, errkit.AddFuncName("repository.(*UserRepositoryImpl).CountByUsername", err)
	}
	return total, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, db *gorm.DB, entity *entity.User, id int64) error {
	err := db.Where("id = ?", id).Take(entity).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName("repository.(*UserRepositoryImpl).FindByID", err)
	}
	return nil
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, db *gorm.DB, entity *entity.User, username string) error {
	err := db.Where("username = ?", username).Take(entity).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName("repository.(*UserRepositoryImpl).FindByUsername", err)
	}
	return nil
}
