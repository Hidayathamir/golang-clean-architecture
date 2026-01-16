package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/CommentRepository.go -pkg=mock . CommentRepository

type CommentRepository interface {
	Create(ctx context.Context, db *gorm.DB, comment *entity.Comment) error
}

var _ CommentRepository = &CommentRepositoryImpl{}

type CommentRepositoryImpl struct {
	Config *viper.Viper
}

func NewCommentRepository(cfg *viper.Viper) *CommentRepositoryImpl {
	return &CommentRepositoryImpl{
		Config: cfg,
	}
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, db *gorm.DB, comment *entity.Comment) error {
	err := db.Create(comment).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
