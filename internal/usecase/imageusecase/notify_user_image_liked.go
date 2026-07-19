package imageusecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func (u *ImageUsecaseImpl) NotifyUserImageLiked(ctx context.Context, req dto.NotifyUserImageLikedRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageLiked")
	}

	image := entity.Image{}

	err = u.ImageRepository.FindByID(ctx, u.DB, &image, req.ImageID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageLiked")
	}

	uploader := entity.User{}

	err = u.UserRepository.FindByID(ctx, u.DB, &uploader, image.UserID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageLiked")
	}

	liker := entity.User{}

	err = u.UserRepository.FindByID(ctx, u.DB, &liker, req.LikerUserID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageLiked")
	}

	event := dto.NotifEvent{
		UserID:  uploader.ID,
		Message: fmt.Sprintf("%s just liked your post", liker.Name),
	}

	err = u.NotifProducer.SendNotif(ctx, &event)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageLiked")
	}

	return nil
}
