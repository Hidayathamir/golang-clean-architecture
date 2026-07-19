package userusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Update(ctx context.Context, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.UserResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Update")
	}

	user := entity.User{}
	err = u.UserRepository.FindByID(ctx, u.DB, &user, req.ID)
	if err != nil {
		return dto.UserResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Update")
	}

	var password string
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return dto.UserResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Update")
		}
		password = string(hashedPassword)
	}

	converter.DtoUpdateUserRequestToEntityUser(req, &user, password)

	err = u.UserRepository.Update(ctx, u.DB, &user)
	if err != nil {
		return dto.UserResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Update")
	}

	err = u.UserCache.Delete(ctx, req.ID)
	logkit.LogIfErr(err)

	res := dto.UserResponse{}
	converter.EntityUserToDtoUserResponse(user, &res)

	return res, nil
}
