package model

import (
	"strconv"
	"time"
)

type UserFollowedEvent struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserFollowedEvent) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}
