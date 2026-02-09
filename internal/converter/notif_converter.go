package converter

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
)

func DtoNotifEventToDtoNotifyRequest(ctx context.Context, event *dto.NotifEvent, req *dto.NotifyRequest) {
	req.UserID = event.UserID
	req.Message = event.Message
}
