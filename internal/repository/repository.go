package repository

import (
	"golang-clean-architecture/pkg/errkit"

	"gorm.io/gorm"
)

//go:generate moq -out=../mock/Repository.go -pkg=mock . Repository

type Repository[T any] interface {
	Create(db *gorm.DB, entity *T) error
	Update(db *gorm.DB, entity *T) error
	Delete(db *gorm.DB, entity *T) error
	CountById(db *gorm.DB, id any) (int64, error)
	FindById(db *gorm.DB, entity *T, id any) error
}

var _ Repository[any] = &RepositoryImpl[any]{}

type RepositoryImpl[T any] struct {
}

func (r *RepositoryImpl[T]) Create(db *gorm.DB, entity *T) error {
	err := db.Create(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *RepositoryImpl[T]) Update(db *gorm.DB, entity *T) error {
	err := db.Save(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *RepositoryImpl[T]) Delete(db *gorm.DB, entity *T) error {
	err := db.Delete(entity).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *RepositoryImpl[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("id = ?", id).Count(&total).Error
	if err != nil {
		return 0, errkit.AddFuncName(err)
	}
	return total, nil
}

func (r *RepositoryImpl[T]) FindById(db *gorm.DB, entity *T, id any) error {
	err := db.Where("id = ?", id).Take(entity).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}
