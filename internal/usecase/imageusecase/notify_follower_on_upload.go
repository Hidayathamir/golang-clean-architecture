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

func (u *ImageUsecaseImpl) NotifyFollowerOnUpload(ctx context.Context, req dto.NotifyFollowerOnUploadRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyFollowerOnUpload")
	}

	user := entity.User{}

	err = u.UserRepository.FindByID(ctx, u.DB, &user, req.UserID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyFollowerOnUpload")
	}

	followList := entity.FollowList{}

	err = u.FollowRepository.FindByFollowingID(ctx, u.DB, &followList, user.ID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyFollowerOnUpload")
	}

	for _, follow := range followList {
		event := dto.NotifEvent{
			UserID:  follow.FollowerID,
			Message: fmt.Sprintf("%s just upload an image", user.Name),
		}

		err = u.NotifProducer.SendNotif(ctx, &event)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyFollowerOnUpload")
		}
	}

	return nil
}
