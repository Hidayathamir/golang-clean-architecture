package model

import "github.com/golang-jwt/jwt/v5"

type RegisterUser2Request struct {
	Email       string `json:"email"        validate:"required,email,max=255"`
	Password    string `json:"password"     validate:"required,min=8,max=100"`
	DisplayName string `json:"display_name" validate:"required,max=100"`
}

type LoginUser2Request struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=100"`
}

type User2Response struct {
	ID          string `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	CreatedAt   int64  `json:"created_at,omitempty"`
	UpdatedAt   int64  `json:"updated_at,omitempty"`
}

type User2TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type User2Auth struct {
	UserID string
}

type VerifyUser2TokenRequest struct {
	Token string `validate:"required"`
}

type GetUser2Request struct {
	ID string `validate:"required"`
}

type User2Claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}
