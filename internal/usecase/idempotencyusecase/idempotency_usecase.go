package idempotencyusecase

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseIdempotency.go -pkg=mock . IdempotencyUsecase

type IdempotencyUsecase interface {
	InsertIfNotExists(ctx context.Context, key, topic string, partition int32, offset int64) (bool, error)
	DeleteOlderThan(ctx context.Context, age time.Duration) (int64, error)
}

var _ IdempotencyUsecase = &IdempotencyUsecaseImpl{}

type IdempotencyUsecaseImpl struct {
	Config                *config.Config
	DB                    *gorm.DB
	IdempotencyRepository repository.IdempotencyRepository
}

func NewIdempotencyUsecase(
	cfg *config.Config,
	db *gorm.DB,
	idempotencyRepository repository.IdempotencyRepository,
) *IdempotencyUsecaseImpl {
	return &IdempotencyUsecaseImpl{
		Config:                cfg,
		DB:                    db,
		IdempotencyRepository: idempotencyRepository,
	}
}
