package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Create(ctx context.Context, req *dto.RegisterUserRequest) (*dto.UserResponse, error) {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	converter.DtoRegisterUserRequestToEntityUser(req, user, string(password))

	if err := u.UserRepository.Create(ctx, u.DB.WithContext(ctx), user); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(dto.UserResponse)
	converter.EntityUserToDtoUserResponse(user, res)

	return res, nil
}
