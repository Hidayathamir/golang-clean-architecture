package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) Login(ctx context.Context, req dto.LoginUserRequest) (dto.UserLoginResponse, error) {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		return dto.UserLoginResponse{}, errkit.AddFuncName(err)
	}

	user := entity.User{}
	err = u.UserRepository.FindByUsername(ctx, u.DB.WithContext(ctx), &user, req.Username)
	if err != nil {
		err = errkit.Unauthorized(err)
		return dto.UserLoginResponse{}, errkit.AddFuncName(err)
	}

	err = user.ValidatePassword(req.Password)
	if err != nil {
		err = errkit.Unauthorized(err)
		return dto.UserLoginResponse{}, errkit.AddFuncName(err)
	}

	token, err := u.signAccessToken(ctx, user.ID)
	if err != nil {
		return dto.UserLoginResponse{}, errkit.AddFuncName(err)
	}

	res := dto.UserLoginResponse{}
	converter.EntityUserToDtoUserLoginResponse(user, &res)
	res.Token = token

	return res, nil
}
