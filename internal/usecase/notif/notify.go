package notif

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *NotifUsecaseImpl) Notify(ctx context.Context, req *model.NotifyRequest) error {
	x.Logger.WithContext(ctx).WithField("req", req).Info("dummy notify")
	return nil
}
