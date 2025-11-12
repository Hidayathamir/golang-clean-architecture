package model

import "strconv"

type UserEvent struct {
	ID        int64  `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	Name      string `json:"name,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}

func (u *UserEvent) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}
