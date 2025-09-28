package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(ctx, tx, user, req.ID); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName(err)
	}

	user.Token = uuid.New().String()
	if err := u.UserRepository.Update(ctx, tx, user); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if _, err := u.SlackClient.IsConnected(ctx); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err := u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToTokenResponse(user), nil
}
