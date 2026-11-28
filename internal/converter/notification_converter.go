package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
)

func EntityNotificationToDtoNotificationResponse(notification entity.Notification, res *dto.NotificationResponse) {
	res.ID = notification.ID
	res.ActorID = notification.ActorID
	res.Type = notification.Type
	res.ImageID = notification.ImageID
	res.IsRead = notification.IsRead
	res.CreatedAt = notification.CreatedAt
}

func EntityNotificationListToDtoNotificationResponseList(notificationList entity.NotificationList, res *dto.NotificationResponseList) {
	for i := range notificationList {
		item := dto.NotificationResponse{}
		EntityNotificationToDtoNotificationResponse(notificationList[i], &item)
		*res = append(*res, item)
	}
}

func DtoNotifyRequestToEntityNotification(req dto.NotifyRequest, notification *entity.Notification) {
	notification.UserID = req.UserID
	notification.ActorID = req.ActorID
	notification.Type = req.Type
	notification.ImageID = req.ImageID
}
