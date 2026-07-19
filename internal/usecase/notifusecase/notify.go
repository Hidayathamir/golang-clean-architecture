package notifusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
)

func (u *NotifUsecaseImpl) Notify(ctx context.Context, req dto.NotifyRequest) error {
	logkit.Logger.WithContext(ctx).WithField("req", req).Info("dummy notify")
	return nil
}
