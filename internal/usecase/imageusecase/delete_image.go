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

	err = u.ImageRepository.DeleteByID(ctx, u.DB, &image)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).DeleteImage")
	}

	return nil
}
