package user

import (
	"context"
	"errors"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	err := u.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	total, err := u.UserRepository.CountById(ctx, u.DB.WithContext(ctx), req.ID)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if total > 0 {
		err = errors.New("user already exists")
		err = errkit.Conflict(err)
		return nil, errkit.AddFuncName(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	user := &entity.User{
		ID:       req.ID,
		Password: string(password),
		Name:     req.Name,
	}

	if err := u.UserRepository.Create(ctx, u.DB.WithContext(ctx), user); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err = u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToResponse(user), nil
}
