package model

import "time"

type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VerifyUserRequest struct {
	Token string `validate:"required"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Name     string `json:"name"     validate:"required,max=100"`
}

type UpdateUserRequest struct {
	ID       int64  `json:"-"        validate:"required"`
	Password string `json:"password" validate:"max=100"`
	Name     string `json:"name"     validate:"max=100"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LogoutUserRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type GetUserRequest struct {
	ID int64 `json:"id" validate:"required"`
}
