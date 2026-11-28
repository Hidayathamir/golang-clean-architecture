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
)

func (u *ImageUsecaseImpl) UpdateImage(ctx context.Context, req dto.UpdateImageRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).UpdateImage")
	}

	userAuth := ctxuserauth.Get(ctx)

	image := entity.Image{}
	err = u.ImageRepository.FindByID(ctx, u.DB, &image, req.ID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).UpdateImage")
	}

	if image.UserID != userAuth.ID {
		err := errors.New("only the owner may edit this image")
		err = errkit.SetCode(err, http.StatusForbidden)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).UpdateImage")
	}

	err = u.ImageRepository.UpdateCaptionByID(ctx, u.DB, req.ID, req.Caption)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).UpdateImage")
	}

	return nil
}
