package image

import (
	"context"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) NotifyFollowerOnUpload(ctx context.Context, req *dto.NotifyFollowerOnUploadRequest) error {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	user := new(entity.User)

	if err := u.UserRepository.FindByID(ctx, u.DB, user, req.UserID); err != nil {
		return errkit.AddFuncName(err)
	}

	followList := new(entity.FollowList)

	if err := u.FollowRepository.FindByFollowingID(ctx, u.DB, followList, user.ID); err != nil {
		return errkit.AddFuncName(err)
	}

	for _, follow := range *followList {
		event := &dto.NotifEvent{
			UserID:  follow.FollowerID,
			Message: fmt.Sprintf("%s just upload an image", user.Name),
		}

		if err := u.NotifProducer.SendNotif(ctx, event); err != nil {
			return errkit.AddFuncName(err)
		}
	}

	return nil
}
