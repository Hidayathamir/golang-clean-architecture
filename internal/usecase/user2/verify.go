package user2

import (
	"context"
	"errors"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/golang-jwt/jwt/v5"
)

func (u *User2UsecaseImpl) VerifyToken(ctx context.Context, req *model.VerifyUser2TokenRequest) (*model.User2Auth, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).VerifyToken", err)
	}

	claims := &model.User2Claims{}
	token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errkit.Unauthorized(errors.New("unexpected signing method"))
		}
		return u.jwtSecret, nil
	})
	if err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).VerifyToken", err)
	}

	if !token.Valid || claims.UserID == "" {
		err = errkit.Unauthorized(errors.New("invalid token"))
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).VerifyToken", err)
	}

	user := new(entity.User2)
	if err := u.User2Repository.FindByID(ctx, u.DB.WithContext(ctx), user, claims.UserID); err != nil {
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).VerifyToken", err)
	}

	return &model.User2Auth{UserID: user.ID}, nil
}
