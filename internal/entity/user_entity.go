package entity

import (
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Username  string         `gorm:"column:username"`
	Password  string         `gorm:"column:password"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (u *User) TableName() string {
	return table.User
}

func (u *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}

type UserList []User
