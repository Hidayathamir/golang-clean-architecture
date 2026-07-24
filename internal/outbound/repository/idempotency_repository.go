package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate moq -out=../../mock/MockRepositoryIdempotency.go -pkg=mock . IdempotencyRepository

type IdempotencyRepository interface {
	InsertIfNotExists(ctx context.Context, db *gorm.DB, key, topic string, partition int32, offset int64) (bool, error)
	DeleteOlderThan(ctx context.Context, db *gorm.DB, age time.Duration) (int64, error)
}

var _ IdempotencyRepository = &IdempotencyRepositoryImpl{}

type IdempotencyRepositoryImpl struct {
	Cfg *config.Config
}

func NewIdempotencyRepository(cfg *config.Config) *IdempotencyRepositoryImpl {
	return &IdempotencyRepositoryImpl{
		Cfg: cfg,
	}
}

func (r *IdempotencyRepositoryImpl) InsertIfNotExists(ctx context.Context, db *gorm.DB, key, topic string, partition int32, offset int64) (bool, error) {
	record := entity.MessageIdempotency{
		IdempotencyKey: key,
		Topic:          topic,
		Partition:      partition,
		RecordOffset:   offset,
	}

	result := db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&record)
	if result.Error != nil {
		return false, errkit.AddFuncName(result.Error, "repository.(*IdempotencyRepositoryImpl).InsertIfNotExists")
	}

	return result.RowsAffected > 0, nil
}

func (r *IdempotencyRepositoryImpl) DeleteOlderThan(ctx context.Context, db *gorm.DB, age time.Duration) (int64, error) {
	cutoff := time.Now().Add(-age)

	result := db.WithContext(ctx).
		Where("processed_at < ?", cutoff).
		Delete(&entity.MessageIdempotency{})
	if result.Error != nil {
		return 0, errkit.AddFuncName(result.Error, "repository.(*IdempotencyRepositoryImpl).DeleteOlderThan")
	}

	return result.RowsAffected, nil
}
