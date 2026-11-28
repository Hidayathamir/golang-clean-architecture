package notifusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *NotifUsecaseImpl) Notify(ctx context.Context, req dto.NotifyRequest) error {
	err := u.createAndPush(ctx, req)
	if err != nil {
		return errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).Notify")
	}
	return nil
}
