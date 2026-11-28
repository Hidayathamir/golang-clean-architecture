package imageusecase

import (
	"context"
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
	notifReq := dto.NotifyRequest{
		UserID:  image.UserID,
		ActorID: req.CommenterUserID,
		Type:    entity.NotificationTypeComment,
		ImageID: req.ImageID,
	}
	err = u.NotifUsecase.Notify(ctx, notifReq)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).NotifyUserImageCommented")
	}
	return nil
}
