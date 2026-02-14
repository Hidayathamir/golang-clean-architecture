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

func (u *UserUsecaseImpl) Update(ctx context.Context, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		return dto.UserResponse{}, errkit.AddFuncName(err)
	}

	user := entity.User{}
	err = u.UserRepository.FindByID(ctx, u.DB.WithContext(ctx), &user, req.ID)
	if err != nil {
		return dto.UserResponse{}, errkit.AddFuncName(err)
	}

	var password string
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return dto.UserResponse{}, errkit.AddFuncName(err)
		}
		password = string(hashedPassword)
	}

	converter.DtoUpdateUserRequestToEntityUser(req, &user, password)

	err = u.UserRepository.Update(ctx, u.DB.WithContext(ctx), &user)
	if err != nil {
		return dto.UserResponse{}, errkit.AddFuncName(err)
	}

	err = u.UserCache.Delete(ctx, req.ID)
	x.LogIfErr(err)

	res := dto.UserResponse{}
	converter.EntityUserToDtoUserResponse(user, &res)

	return res, nil
}
