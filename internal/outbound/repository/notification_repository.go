package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockRepositoryNotification.go -pkg=mock . NotificationRepository
type NotificationRepository interface {
	Create(ctx context.Context, db *gorm.DB, notification *entity.Notification) error
	FindByID(ctx context.Context, db *gorm.DB, notification *entity.Notification, id int64) error
	FindByUserID(ctx context.Context, db *gorm.DB, notificationList *entity.NotificationList, userID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error
	MarkReadByID(ctx context.Context, db *gorm.DB, userID int64, notificationID int64) (int64, error)
	MarkAllReadByUserID(ctx context.Context, db *gorm.DB, userID int64) (int64, error)
	CountUnreadByUserID(ctx context.Context, db *gorm.DB, userID int64) (int64, error)
}

var _ NotificationRepository = &NotificationRepositoryImpl{}

type NotificationRepositoryImpl struct {
	Cfg *config.Config
}

func NewNotificationRepository(cfg *config.Config) *NotificationRepositoryImpl {
	return &NotificationRepositoryImpl{
		Cfg: cfg,
	}
}

func (r *NotificationRepositoryImpl) Create(ctx context.Context, db *gorm.DB, notification *entity.Notification) error {
	err := db.WithContext(ctx).Create(notification).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*NotificationRepositoryImpl).Create")
	}
	return nil
}

func (r *NotificationRepositoryImpl) FindByID(ctx context.Context, db *gorm.DB, notification *entity.Notification, id int64) error {
	err := db.WithContext(ctx).Where(column.ID.Eq(id)).Take(notification).Error
	if err != nil {
		err = errkit.SetCode(err, http.StatusNotFound)
		return errkit.AddFuncName(err, "repository.(*NotificationRepositoryImpl).FindByID")
	}
	return nil
}

func (r *NotificationRepositoryImpl) FindByUserID(ctx context.Context, db *gorm.DB, notificationList *entity.NotificationList, userID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	query := db.WithContext(ctx).
		Where(column.UserID.Eq(userID)).
		Order(column.CreatedAt.Desc()).
		Order(column.ID.Desc()).
		Limit(limit)
	if cursorID > 0 {
		whereKeyset, argCreatedAt, argID := column.KeysetBefore(cursorCreatedAt, cursorID)
		query = query.Where(whereKeyset, argCreatedAt, argID)
	}
	err := query.Find(notificationList).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*NotificationRepositoryImpl).FindByUserID")
	}
	return nil
}

func (r *NotificationRepositoryImpl) MarkReadByID(ctx context.Context, db *gorm.DB, userID int64, notificationID int64) (int64, error) {
	whereUser, argUser := column.UserID.Eq(userID)
	whereID, argID := column.ID.Eq(notificationID)
	result := db.WithContext(ctx).
		Model(&entity.Notification{}).
		Where(whereUser, argUser).
		Where(whereID, argID).
		Update(column.IsRead.Str(), true)
	if result.Error != nil {
		return 0, errkit.AddFuncName(result.Error, "repository.(*NotificationRepositoryImpl).MarkReadByID")
	}
	return result.RowsAffected, nil
}

func (r *NotificationRepositoryImpl) MarkAllReadByUserID(ctx context.Context, db *gorm.DB, userID int64) (int64, error) {
	whereUser, argUser := column.UserID.Eq(userID)
	whereUnread, argUnread := column.IsRead.Eq(false)
	result := db.WithContext(ctx).
		Model(&entity.Notification{}).
		Where(whereUser, argUser).
		Where(whereUnread, argUnread).
		Update(column.IsRead.Str(), true)
	if result.Error != nil {
		return 0, errkit.AddFuncName(result.Error, "repository.(*NotificationRepositoryImpl).MarkAllReadByUserID")
	}
	return result.RowsAffected, nil
}

func (r *NotificationRepositoryImpl) CountUnreadByUserID(ctx context.Context, db *gorm.DB, userID int64) (int64, error) {
	var total int64
	whereUser, argUser := column.UserID.Eq(userID)
	whereUnread, argUnread := column.IsRead.Eq(false)
	err := db.WithContext(ctx).
		Model(&entity.Notification{}).
		Where(whereUser, argUser).
		Where(whereUnread, argUnread).
		Count(&total).Error
	if err != nil {
		return 0, errkit.AddFuncName(err, "repository.(*NotificationRepositoryImpl).CountUnreadByUserID")
	}
	return total, nil
}
