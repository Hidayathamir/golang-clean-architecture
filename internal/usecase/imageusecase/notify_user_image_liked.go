package imageusecase

import (
	"context"
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
	notifReq := dto.NotifyRequest{
		UserID:  image.UserID,
		ActorID: req.LikerUserID,
		Type:    entity.NotificationTypeLike,
		ImageID: req.ImageID,
	}
	err = u.NotifUsecase.Notify(ctx, notifReq)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageLiked")
	}
	return nil
}
