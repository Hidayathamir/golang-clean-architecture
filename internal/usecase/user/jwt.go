package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *UserUsecaseImpl) signAccessToken(ctx context.Context, userID int64) (string, error) {
	secret := u.Config.GetAuthJWTSecret()
	if secret == "" {
		err := fmt.Errorf("jwt secret is not configured")
		err = errkit.InternalServerError(err)
		return "", errkit.AddFuncName(err)
	}

	expireSeconds := u.Config.GetAuthJWTExpireSeconds()
	if expireSeconds <= 0 {
		err := fmt.Errorf("jwt expire seconds must be greater than zero")
		err = errkit.InternalServerError(err)
		return "", errkit.AddFuncName(err)
	}

	issuer := u.Config.GetAuthJWTIssuer()
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expireSeconds) * time.Second)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		err = errkit.InternalServerError(err)
		return "", errkit.AddFuncName(err)
	}

	return tokenString, nil
}

func (u *UserUsecaseImpl) parseAccessToken(ctx context.Context, tokenString string) (int64, error) {
	if tokenString == "" {
		err := fmt.Errorf("token is empty")
		err = errkit.Unauthorized(err)
		return 0, errkit.AddFuncName(err)
	}

	secret := u.Config.GetAuthJWTSecret()
	if secret == "" {
		err := fmt.Errorf("jwt secret is not configured")
		err = errkit.InternalServerError(err)
		return 0, errkit.AddFuncName(err)
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		err = errkit.Unauthorized(err)
		return 0, errkit.AddFuncName(err)
	}

	if !token.Valid {
		err := fmt.Errorf("token is invalid")
		err = errkit.Unauthorized(err)
		return 0, errkit.AddFuncName(err)
	}

	if claims.Subject == "" {
		err := fmt.Errorf("token subject is empty")
		err = errkit.Unauthorized(err)
		return 0, errkit.AddFuncName(err)
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		err = errkit.Unauthorized(err)
		return 0, errkit.AddFuncName(err)
	}

	return userID, nil
}
