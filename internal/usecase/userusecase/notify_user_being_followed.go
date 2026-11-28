package userusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func (u *UserUsecaseImpl) NotifyUserBeingFollowed(ctx context.Context, req dto.NotifyUserBeingFollowedRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).NotifyUserBeingFollowed")
	}
	notifReq := dto.NotifyRequest{
		UserID:  req.FollowingID,
		ActorID: req.FollowerID,
		Type:    entity.NotificationTypeFollow,
	}
	err = u.NotifUsecase.Notify(ctx, notifReq)
	if err != nil {
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).NotifyUserBeingFollowed")
	}
	return nil
}
