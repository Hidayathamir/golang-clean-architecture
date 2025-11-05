package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/ContactRepository.go -pkg=mock . ContactRepository

type ContactRepository interface {
	FindByIDAndUserID(ctx context.Context, db *gorm.DB, contact *entity.Contact, id string, userID string) error
	Search(ctx context.Context, db *gorm.DB, req *model.SearchContactRequest) (entity.ContactList, int64, error)
	Create(ctx context.Context, db *gorm.DB, entity *entity.Contact) error
	Update(ctx context.Context, db *gorm.DB, entity *entity.Contact) error
	Delete(ctx context.Context, db *gorm.DB, entity *entity.Contact) error
}

var _ ContactRepository = &ContactRepositoryImpl{}

type ContactRepositoryImpl struct {
	Config *viper.Viper
}

func NewContactRepository(cfg *viper.Viper) *ContactRepositoryImpl {
	return &ContactRepositoryImpl{
		Config: cfg,
	}
}

func (r *ContactRepositoryImpl) FindByIDAndUserID(ctx context.Context, db *gorm.DB, contact *entity.Contact, id string, userID string) error {
	err := db.Where("id = ? AND user_id = ?", id, userID).Take(contact).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName("repository.(*ContactRepositoryImpl).FindByIDAndUserID", err)
	}
	return nil
}

func (r *ContactRepositoryImpl) Search(ctx context.Context, db *gorm.DB, req *model.SearchContactRequest) (entity.ContactList, int64, error) {
	var contacts entity.ContactList
	if err := db.Scopes(r.filterContact(req)).Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&contacts).Error; err != nil {
		return nil, 0, errkit.AddFuncName("repository.(*ContactRepositoryImpl).Search", err)
	}

	var total int64 = 0
	if err := db.Model(&entity.Contact{}).Scopes(r.filterContact(req)).Count(&total).Error; err != nil {
		return nil, 0, errkit.AddFuncName("repository.(*ContactRepositoryImpl).Search", err)
	}

	return contacts, total, nil
}

func (r *ContactRepositoryImpl) filterContact(req *model.SearchContactRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("user_id = ?", req.UserID)

		if name := req.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("first_name ILIKE ? OR last_name ILIKE ?", name, name)
		}

		if phone := req.Phone; phone != "" {
			phone = "%" + phone + "%"
			tx = tx.Where("phone ILIKE ?", phone)
		}

		if email := req.Email; email != "" {
			email = "%" + email + "%"
			tx = tx.Where("email ILIKE ?", email)
		}

		return tx
	}
}

func (r *ContactRepositoryImpl) Create(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := db.Create(entity).Error
	if err != nil {
		return errkit.AddFuncName("repository.(*ContactRepositoryImpl).Create", err)
	}
	return nil
}

func (r *ContactRepositoryImpl) Update(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := db.Save(entity).Error
	if err != nil {
		return errkit.AddFuncName("repository.(*ContactRepositoryImpl).Update", err)
	}
	return nil
}

func (r *ContactRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := db.Delete(entity).Error
	if err != nil {
		return errkit.AddFuncName("repository.(*ContactRepositoryImpl).Delete", err)
	}
	return nil
}
