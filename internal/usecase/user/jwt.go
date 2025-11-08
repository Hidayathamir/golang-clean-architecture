package user

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *UserUsecaseImpl) signAccessToken(ctx context.Context, userID string) (string, error) {
	secret := u.Config.GetString(configkey.AuthJWTSecret)
	if secret == "" {
		err := fmt.Errorf("jwt secret is not configured")
		err = errkit.InternalServerError(err)
		return "", errkit.AddFuncName("user.(*UserUsecaseImpl).signAccessToken", err)
	}

	expireSeconds := u.Config.GetInt(configkey.AuthJWTExpireSeconds)
	if expireSeconds <= 0 {
		err := fmt.Errorf("jwt expire seconds must be greater than zero")
		err = errkit.InternalServerError(err)
		return "", errkit.AddFuncName("user.(*UserUsecaseImpl).signAccessToken", err)
	}

	issuer := u.Config.GetString(configkey.AuthJWTIssuer)
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expireSeconds) * time.Second)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		err = errkit.InternalServerError(err)
		return "", errkit.AddFuncName("user.(*UserUsecaseImpl).signAccessToken", err)
	}

	return tokenString, nil
}

func (u *UserUsecaseImpl) parseAccessToken(ctx context.Context, tokenString string) (string, error) {
	if tokenString == "" {
		err := fmt.Errorf("token is empty")
		err = errkit.Unauthorized(err)
		return "", errkit.AddFuncName("user.(*UserUsecaseImpl).parseAccessToken", err)
	}

	secret := u.Config.GetString(configkey.AuthJWTSecret)
	if secret == "" {
		err := fmt.Errorf("jwt secret is not configured")
		err = errkit.InternalServerError(err)
		return "", errkit.AddFuncName("user.(*UserUsecaseImpl).parseAccessToken", err)
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
		return "", errkit.AddFuncName("user.(*UserUsecaseImpl).parseAccessToken", err)
	}

	if !token.Valid {
		err := fmt.Errorf("token is invalid")
		err = errkit.Unauthorized(err)
		return "", errkit.AddFuncName("user.(*UserUsecaseImpl).parseAccessToken", err)
	}

	if claims.Subject == "" {
		err := fmt.Errorf("token subject is empty")
		err = errkit.Unauthorized(err)
		return "", errkit.AddFuncName("user.(*UserUsecaseImpl).parseAccessToken", err)
	}

	return claims.Subject, nil
}
