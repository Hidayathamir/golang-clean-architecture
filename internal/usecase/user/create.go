package user

import (
	"context"
	"errors"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/pkg/errkit"

	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	total, err := u.UserRepository.CountById(ctx, tx, req.ID)
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

	if err := u.UserRepository.Create(ctx, tx, user); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err = u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToResponse(user), nil
}
