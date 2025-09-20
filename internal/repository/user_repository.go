package repository

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/pkg/errkit"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/UserRepository.go -pkg=mock . UserRepository

type UserRepository interface {
	FindByToken(db *gorm.DB, user *entity.User, token string) error
	Create(db *gorm.DB, entity *entity.User) error
	Update(db *gorm.DB, entity *entity.User) error
	CountById(db *gorm.DB, id string) (int64, error)
	FindById(db *gorm.DB, entity *entity.User, id string) error
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

func (r *UserRepositoryImpl) FindByToken(db *gorm.DB, user *entity.User, token string) error {
	err := db.Where("token = ?", token).First(user).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) Create(db *gorm.DB, entity *entity.User) error {
	err := db.Create(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) Update(db *gorm.DB, entity *entity.User) error {
	err := db.Save(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(db *gorm.DB, entity *entity.User) error {
	err := db.Delete(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *UserRepositoryImpl) CountById(db *gorm.DB, id string) (int64, error) {
	var total int64
	err := db.Model(new(entity.User)).Where("id = ?", id).Count(&total).Error
	if err != nil {
		return 0, errkit.AddFuncName(err)
	}
	return total, nil
}

func (r *UserRepositoryImpl) FindById(db *gorm.DB, entity *entity.User, id string) error {
	err := db.Where("id = ?", id).Take(entity).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}
