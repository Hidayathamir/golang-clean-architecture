package repository

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/errkit"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/ContactRepository.go -pkg=mock . ContactRepository

type ContactRepository interface {
	Repository[entity.Contact]
	FindByIdAndUserId(db *gorm.DB, contact *entity.Contact, id string, userId string) error
	Search(db *gorm.DB, request *model.SearchContactRequest) ([]entity.Contact, int64, error)
}

var _ ContactRepository = &ContactRepositoryImpl{}

type ContactRepositoryImpl struct {
	RepositoryImpl[entity.Contact]
	Log *logrus.Logger
}

func NewContactRepository(log *logrus.Logger) *ContactRepositoryImpl {
	return &ContactRepositoryImpl{
		Log: log,
	}
}

func (r *ContactRepositoryImpl) FindByIdAndUserId(db *gorm.DB, contact *entity.Contact, id string, userId string) error {
	err := db.Where("id = ? AND user_id = ?", id, userId).Take(contact).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *ContactRepositoryImpl) Search(db *gorm.DB, request *model.SearchContactRequest) ([]entity.Contact, int64, error) {
	var contacts []entity.Contact
	if err := db.Scopes(r.filterContact(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&contacts).Error; err != nil {
		return nil, 0, errkit.AddFuncName(err)
	}

	var total int64 = 0
	if err := db.Model(&entity.Contact{}).Scopes(r.filterContact(request)).Count(&total).Error; err != nil {
		return nil, 0, errkit.AddFuncName(err)
	}

	return contacts, total, nil
}

func (r *ContactRepositoryImpl) filterContact(request *model.SearchContactRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("user_id = ?", request.UserId)

		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("first_name LIKE ? OR last_name LIKE ?", name, name)
		}

		if phone := request.Phone; phone != "" {
			phone = "%" + phone + "%"
			tx = tx.Where("phone LIKE ?", phone)
		}

		if email := request.Email; email != "" {
			email = "%" + email + "%"
			tx = tx.Where("email LIKE ?", email)
		}

		return tx
	}
}
