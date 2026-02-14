package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) Current(ctx context.Context, req dto.GetUserRequest) (dto.UserResponse, error) {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		return dto.UserResponse{}, errkit.AddFuncName(err)
	}

	cachedUser, err := u.UserCache.Get(ctx, req.ID)
	if err == nil && cachedUser != nil {
		res := dto.UserResponse{}
		converter.EntityUserToDtoUserResponse(*cachedUser, &res)
		return res, nil
	}

	user := entity.User{}
	err = u.UserRepository.FindByID(ctx, u.DB.WithContext(ctx), &user, req.ID)
	if err != nil {
		return dto.UserResponse{}, errkit.AddFuncName(err)
	}

	err = u.UserCache.Set(ctx, &user)
	x.LogIfErr(err)

	res := dto.UserResponse{}
	converter.EntityUserToDtoUserResponse(user, &res)

	return res, nil
}
