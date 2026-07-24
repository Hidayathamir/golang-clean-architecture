package idempotencyusecase

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *IdempotencyUsecaseImpl) DeleteOlderThan(ctx context.Context, age time.Duration) (int64, error) {
	deleted, err := u.IdempotencyRepository.DeleteOlderThan(ctx, u.DB, age)
	if err != nil {
		return 0, errkit.AddFuncName(err, "idempotencyusecase.(*IdempotencyUsecaseImpl).DeleteOlderThan")
	}
	return deleted, nil
}
