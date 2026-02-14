package image

import (
	"context"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) NotifyUserImageCommented(ctx context.Context, req dto.NotifyUserImageCommentedRequest) error {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	image := entity.Image{}

	err = u.ImageRepository.FindByID(ctx, u.DB, &image, req.ImageID)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	uploader := entity.User{}

	err = u.UserRepository.FindByID(ctx, u.DB, &uploader, image.UserID)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	commenter := entity.User{}

	err = u.UserRepository.FindByID(ctx, u.DB, &commenter, req.CommenterUserID)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	event := dto.NotifEvent{
		UserID:  uploader.ID,
		Message: fmt.Sprintf("%s just comment on your post", commenter.Name),
	}

	err = u.NotifProducer.SendNotif(ctx, &event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
