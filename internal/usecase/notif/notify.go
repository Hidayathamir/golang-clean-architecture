package notif

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *NotifUsecaseImpl) Notify(ctx context.Context, req *dto.NotifyRequest) error {
	x.Logger.WithContext(ctx).WithField("req", req).Info("dummy notify")
	return nil
}
