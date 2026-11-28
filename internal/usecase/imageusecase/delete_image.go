package imageusecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"gorm.io/gorm"
)

func (u *ImageUsecaseImpl) DeleteImage(ctx context.Context, req dto.DeleteImageRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImage")
	}

	userAuth := ctxuserauth.Get(ctx)

	image := entity.Image{}
	err = u.ImageRepository.FindByID(ctx, u.DB, &image, req.ID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImage")
	}

	if image.UserID != userAuth.ID {
		err := errors.New("only the owner may delete this image")
		err = errkit.SetCode(err, http.StatusForbidden)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImage")
	}

	// The row is soft-deleted first and the S3 object is cleaned up asynchronously
	// via the outbox, so a storage failure never blocks the delete or leaves a
	// visible row pointing at a missing object.
	err = u.DB.Transaction(func(tx *gorm.DB) error {
		err := u.ImageRepository.DeleteByID(ctx, tx, &image)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImage")
		}

		event := dto.ImageDeletedEvent{}
		converter.EntityImageToDtoImageDeletedEvent(image, &event)

		err = u.ImageProducer.SendImageDeleted(ctx, tx, &event)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImage")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
