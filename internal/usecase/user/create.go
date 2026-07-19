package user

import (
	"net/http" 
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Create(ctx context.Context, req dto.RegisterUserRequest) (dto.UserResponse, error) {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.UserResponse{}, errkit.AddFuncName(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, errkit.AddFuncName(err)
	}

	user := entity.User{}
	converter.DtoRegisterUserRequestToEntityUser(req, &user, string(password))

	err = u.UserRepository.Create(ctx, u.DB, &user)
	if err != nil {
		return dto.UserResponse{}, errkit.AddFuncName(err)
	}

	res := dto.UserResponse{}
	converter.EntityUserToDtoUserResponse(user, &res)

	return res, nil
}
