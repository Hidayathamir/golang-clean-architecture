package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserLoginResponse, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByUsername(ctx, u.DB.WithContext(ctx), user, req.Username); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName(err)
	}

	if err := user.ValidatePassword(req.Password); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName(err)
	}

	token, err := u.signAccessToken(ctx, user.ID)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(model.UserLoginResponse)
	converter.EntityUserToModelUserLoginResponse(user, res)
	res.Token = token

	return res, nil
}
