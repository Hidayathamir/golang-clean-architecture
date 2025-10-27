package user2

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (u *User2UsecaseImpl) Login(ctx context.Context, req *model.LoginUser2Request) (*model.User2TokenResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Login", err)
	}

	user := new(entity.User2)
	if err := u.User2Repository.FindByEmail(ctx, u.DB.WithContext(ctx), user, req.Email); err != nil {
		var httpErr *errkit.HTTPError
		if errors.As(err, &httpErr) && httpErr.HTTPCode == http.StatusNotFound {
			err = errkit.Unauthorized(errors.New("invalid credentials"))
		}
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Login", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Login", err)
	}

	now := time.Now().UTC()
	expiresAt := now.Add(u.jwtTTL)

	claims := model.User2Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Login", err)
	}

	return &model.User2TokenResponse{
		AccessToken: signed,
		ExpiresAt:   expiresAt.Unix(),
	}, nil
}
