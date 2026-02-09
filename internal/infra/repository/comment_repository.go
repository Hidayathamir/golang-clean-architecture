package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockRepositoryComment.go -pkg=mock . CommentRepository

type CommentRepository interface {
	Create(ctx context.Context, db *gorm.DB, comment *entity.Comment) error
	FindByImageID(ctx context.Context, db *gorm.DB, commentList *entity.CommentList, imageID int64) error
}

var _ CommentRepository = &CommentRepositoryImpl{}

type CommentRepositoryImpl struct {
	Cfg *config.Config
}

func NewCommentRepository(cfg *config.Config) *CommentRepositoryImpl {
	return &CommentRepositoryImpl{
		Cfg: cfg,
	}
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, db *gorm.DB, comment *entity.Comment) error {
	err := db.WithContext(ctx).Create(comment).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *CommentRepositoryImpl) FindByImageID(ctx context.Context, db *gorm.DB, commentList *entity.CommentList, imageID int64) error {
	err := db.WithContext(ctx).Where(column.ImageID.Eq(imageID)).Find(commentList).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
