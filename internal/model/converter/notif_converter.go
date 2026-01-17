package converter

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
)

func ModelNotifEventToModelNotifyRequest(ctx context.Context, event *model.NotifEvent, req *model.NotifyRequest) {
	req.UserID = event.UserID
	req.Message = event.Message
}
