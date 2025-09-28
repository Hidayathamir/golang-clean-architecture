package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(ctx, tx, user, req.ID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errkit.AddFuncName(err)
		}
		user.Password = string(password)
	}

	if err := u.UserRepository.Update(ctx, tx, user); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err := u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToResponse(user), nil
}
