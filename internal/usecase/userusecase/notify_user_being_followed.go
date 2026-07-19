package userusecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) NotifyUserBeingFollowed(ctx context.Context, req dto.NotifyUserBeingFollowedRequest) error {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).NotifyUserBeingFollowed")
	}

	followerUser := entity.User{}

	err = u.UserRepository.FindByID(ctx, u.DB, &followerUser, req.FollowerID)
	if err != nil {
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).NotifyUserBeingFollowed")
	}

	event := dto.NotifEvent{
		UserID:  req.FollowingID,
		Message: fmt.Sprintf("%s just follow you", followerUser.Name),
	}

	err = u.NotifProducer.SendNotif(ctx, &event)
	if err != nil {
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).NotifyUserBeingFollowed")
	}

	return nil
}
