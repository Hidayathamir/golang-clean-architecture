package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) Current(ctx context.Context, req *dto.GetUserRequest) (*dto.UserResponse, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	cachedUser, err := u.UserCache.Get(ctx, req.ID)
	if err == nil && cachedUser != nil {
		res := new(dto.UserResponse)
		converter.EntityUserToDtoUserResponse(cachedUser, res)
		return res, nil
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByID(ctx, u.DB.WithContext(ctx), user, req.ID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	err = u.UserCache.Set(ctx, user)
	x.LogIfErr(err)

	res := new(dto.UserResponse)
	converter.EntityUserToDtoUserResponse(user, res)

	return res, nil
}
