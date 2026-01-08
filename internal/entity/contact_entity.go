package entity

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	FirstName string         `gorm:"column:first_name"`
	LastName  string         `gorm:"column:last_name"`
	Email     string         `gorm:"column:email"`
	Phone     string         `gorm:"column:phone"`
	UserID    int64          `gorm:"column:user_id"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User      User        `gorm:"foreignKey:user_id;references:id"`
	Addresses AddressList `gorm:"foreignKey:contact_id;references:id"`
}

func (c *Contact) TableName() string {
	return "contacts"
}

type ContactList []Contact
