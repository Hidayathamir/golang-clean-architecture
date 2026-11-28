package userusecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"gorm.io/gorm"
)

func (u *UserUsecaseImpl) Unfollow(ctx context.Context, req dto.UnfollowUserRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Unfollow")
	}

	userAuth := ctxuserauth.Get(ctx)

	event := dto.UserFollowedEvent{
		FollowerID:  userAuth.ID,
		FollowingID: req.FollowingID,
		Count:       -1, // decrement follower/following stats downstream
	}

	err = u.DB.Transaction(func(tx *gorm.DB) error {
		deleted, err := u.FollowRepository.DeleteByFollowerIDAndFollowingID(ctx, tx, userAuth.ID, req.FollowingID)
		if err != nil {
			return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Unfollow")
		}
		if deleted == 0 {
			err := errors.New("follow not found")
			err = errkit.SetCode(err, http.StatusNotFound)
			return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Unfollow")
		}

		err = u.UserProducer.SendUserFollowed(ctx, tx, &event)
		if err != nil {
			return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Unfollow")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
