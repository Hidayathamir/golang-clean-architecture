package userusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func (u *UserUsecaseImpl) Login(ctx context.Context, req dto.LoginUserRequest) (dto.UserLoginResponse, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.UserLoginResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Login")
	}

	user := entity.User{}
	err = u.UserRepository.FindByUsername(ctx, u.DB, &user, req.Username)
	if err != nil {
		err = errkit.SetCode(err, http.StatusUnauthorized)
		return dto.UserLoginResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Login")
	}

	err = user.ValidatePassword(req.Password)
	if err != nil {
		err = errkit.SetCode(err, http.StatusUnauthorized)
		return dto.UserLoginResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Login")
	}

	token, err := u.signAccessToken(ctx, user.ID)
	if err != nil {
		return dto.UserLoginResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Login")
	}

	res := dto.UserLoginResponse{}
	converter.EntityUserToDtoUserLoginResponse(user, &res)
	res.Token = token

	return res, nil
}
