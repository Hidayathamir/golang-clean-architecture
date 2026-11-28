package imageusecase

import (
	"context"
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
	followList := entity.FollowList{}
	err = u.FollowRepository.FindByFollowingID(ctx, u.DB, &followList, req.UserID)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyFollowerOnUpload")
	}
	for i := range followList {
		notifReq := dto.NotifyRequest{
			UserID:  followList[i].FollowerID,
			ActorID: req.UserID,
			Type:    entity.NotificationTypeNewPost,
		}
		err = u.NotifUsecase.Notify(ctx, notifReq)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyFollowerOnUpload")
		}
	}
	return nil
}
