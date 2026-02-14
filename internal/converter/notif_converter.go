package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
)

func DtoNotifEventToDtoNotifyRequest(event dto.NotifEvent, req *dto.NotifyRequest) {
	req.UserID = event.UserID
	req.Message = event.Message
}
