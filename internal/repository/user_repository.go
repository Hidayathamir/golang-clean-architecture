package repository

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/pkg/errkit"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/UserRepository.go -pkg=mock . UserRepository

type UserRepository interface {
	Repository[entity.User]
	FindByToken(db *gorm.DB, user *entity.User, token string) error
}

var _ UserRepository = &UserRepositoryImpl{}

type UserRepositoryImpl struct {
	RepositoryImpl[entity.User]
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
