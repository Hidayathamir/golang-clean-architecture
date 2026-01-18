package image

import (
	"context"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) NotifyUserImageCommented(ctx context.Context, req *model.NotifyUserImageCommentedRequest) error {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	image := new(entity.Image)

	if err := u.ImageRepository.FindByID(ctx, u.DB, image, req.ImageID); err != nil {
		return errkit.AddFuncName(err)
	}

	uploader := new(entity.User)

	if err := u.UserRepository.FindByID(ctx, u.DB, uploader, image.UserID); err != nil {
		return errkit.AddFuncName(err)
	}

	commenter := new(entity.User)

	if err := u.UserRepository.FindByID(ctx, u.DB, commenter, req.CommenterUserID); err != nil {
		return errkit.AddFuncName(err)
	}

	event := &model.NotifEvent{
		UserID:  uploader.ID,
		Message: fmt.Sprintf("%s just comment on your post", commenter.Name),
	}

	if err := u.NotifProducer.SendNotif(ctx, event); err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
