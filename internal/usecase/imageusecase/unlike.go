package imageusecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"gorm.io/gorm"
)

func (u *ImageUsecaseImpl) Unlike(ctx context.Context, req dto.UnlikeImageRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Unlike")
	}

	userAuth := ctxuserauth.Get(ctx)

	event := dto.ImageLikedEvent{
		UserID:  userAuth.ID,
		ImageID: req.ImageID,
		Count:   -1, // decrement the image's like count downstream
	}

	err = u.DB.Transaction(func(tx *gorm.DB) error {
		deleted, err := u.LikeRepository.DeleteByUserIDAndImageID(ctx, tx, userAuth.ID, req.ImageID)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Unlike")
		}
		if deleted == 0 {
			err := errors.New("like not found")
			err = errkit.SetCode(err, http.StatusNotFound)
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Unlike")
		}

		err = u.ImageProducer.SendImageLiked(ctx, tx, &event)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Unlike")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
