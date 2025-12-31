package entity

import "time"

type Address struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	ContactID  int64     `gorm:"column:contact_id"`
	Street     string    `gorm:"column:street"`
	City       string    `gorm:"column:city"`
	Province   string    `gorm:"column:province"`
	PostalCode string    `gorm:"column:postal_code"`
	Country    string    `gorm:"column:country"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`

	Contact Contact `gorm:"foreignKey:contact_id;references:id"`
}

func (a *Address) TableName() string {
	return "addresses"
}

type AddressList []Address
