package model

import (
	"strconv"
	"time"
)

type UserEvent struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserEvent) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}
