package idempotencyusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *IdempotencyUsecaseImpl) InsertIfNotExists(ctx context.Context, key, topic string, partition int32, offset int64) (bool, error) {
	isNew, err := u.IdempotencyRepository.InsertIfNotExists(ctx, u.DB, key, topic, partition, offset)
	if err != nil {
		return false, errkit.AddFuncName(err, "idempotencyusecase.(*IdempotencyUsecaseImpl).InsertIfNotExists")
	}
	return isNew, nil
}
