package model

import "time"

type ContactResponse struct {
	ID        int64               `json:"id"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Email     string              `json:"email"`
	Phone     string              `json:"phone"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Addresses AddressResponseList `json:"addresses"`
}

type ContactResponseList []ContactResponse

type CreateContactRequest struct {
	UserID    int64  `json:"-"          validate:"required"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name"  validate:"max=100"`
	Email     string `json:"email"      validate:"max=200,email"`
	Phone     string `json:"phone"      validate:"max=20"`
}

type UpdateContactRequest struct {
	UserID    int64  `json:"-"          validate:"required"`
	ID        int64  `json:"-"          validate:"required"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name"  validate:"max=100"`
	Email     string `json:"email"      validate:"max=200,email"`
	Phone     string `json:"phone"      validate:"max=20"`
}

type SearchContactRequest struct {
	UserID int64  `json:"-"     validate:"required"`
	Name   string `json:"name"  validate:"max=100"`
	Email  string `json:"email" validate:"max=200"`
	Phone  string `json:"phone" validate:"max=20"`
	Page   int    `json:"page"  validate:"min=1"`
	Size   int    `json:"size"  validate:"min=1,max=100"`
}

type GetContactRequest struct {
	UserID int64 `json:"-" validate:"required"`
	ID     int64 `json:"-" validate:"required"`
}

type DeleteContactRequest struct {
	UserID int64 `json:"-" validate:"required"`
	ID     int64 `json:"-" validate:"required"`
}
