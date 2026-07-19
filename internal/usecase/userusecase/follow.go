package userusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"gorm.io/gorm"
)

func (u *UserUsecaseImpl) Follow(ctx context.Context, req dto.FollowUserRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Follow")
	}

	follow := entity.Follow{}
	converter.DtoFollowUserRequestToEntityFollow(ctx, req, &follow)

	event := dto.UserFollowedEvent{}
	converter.EntityFollowToDtoUserFollowedEvent(follow, &event)

	err = u.DB.Transaction(func(tx *gorm.DB) error {
		err := u.FollowRepository.Create(ctx, tx, &follow)
		if err != nil {
			return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Follow")
		}

		err = u.UserProducer.SendUserFollowed(ctx, tx, &event)
		if err != nil {
			return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Follow")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
