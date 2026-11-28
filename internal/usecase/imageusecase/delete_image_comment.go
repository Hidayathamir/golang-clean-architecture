package imageusecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"gorm.io/gorm"
)

func (u *ImageUsecaseImpl) DeleteImageComment(ctx context.Context, req dto.DeleteImageCommentRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImageComment")
	}

	userAuth := ctxuserauth.Get(ctx)

	comment := entity.Comment{}
	err = u.CommentRepository.FindByID(ctx, u.DB, &comment, req.CommentID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImageComment")
	}

	if comment.ImageID != req.ImageID {
		err := errors.New("comment does not belong to this image")
		err = errkit.SetCode(err, http.StatusNotFound)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImageComment")
	}

	image := entity.Image{}
	err = u.ImageRepository.FindByID(ctx, u.DB, &image, req.ImageID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImageComment")
	}

	if comment.UserID != userAuth.ID && image.UserID != userAuth.ID {
		err := errors.New("only the comment owner or image owner may delete this comment")
		err = errkit.SetCode(err, http.StatusForbidden)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImageComment")
	}

	event := dto.ImageCommentedEvent{
		UserID:  userAuth.ID,
		ImageID: req.ImageID,
		Count:   -1, // decrement the image's comment count downstream
	}

	err = u.DB.Transaction(func(tx *gorm.DB) error {
		err := u.CommentRepository.DeleteByID(ctx, tx, &comment)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImageComment")
		}

		err = u.ImageProducer.SendImageCommented(ctx, tx, &event)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImageComment")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
