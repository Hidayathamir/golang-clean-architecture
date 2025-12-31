package model

import (
	"strconv"
	"time"
)

type UserEvent struct {
	ID        int64     `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u *UserEvent) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}
