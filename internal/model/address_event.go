package model

import "time"

type AddressEvent struct {
	ID         int64     `json:"id"`
	ContactID  int64     `json:"contact_id"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	Province   string    `json:"province"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (a *AddressEvent) GetID() int64 {
	return a.ID
}
