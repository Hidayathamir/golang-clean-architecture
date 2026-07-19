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

func (u *ImageUsecaseImpl) NotifyUserImageCommented(ctx context.Context, req dto.NotifyUserImageCommentedRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageCommented")
	}

	image := entity.Image{}

	err = u.ImageRepository.FindByID(ctx, u.DB, &image, req.ImageID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageCommented")
	}

	uploader := entity.User{}

	err = u.UserRepository.FindByID(ctx, u.DB, &uploader, image.UserID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageCommented")
	}

	commenter := entity.User{}

	err = u.UserRepository.FindByID(ctx, u.DB, &commenter, req.CommenterUserID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageCommented")
	}

	event := dto.NotifEvent{
		UserID:  uploader.ID,
		Message: fmt.Sprintf("%s just comment on your post", commenter.Name),
	}

	err = u.NotifProducer.SendNotif(ctx, u.DB, &event)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageCommented")
	}

	return nil
}
