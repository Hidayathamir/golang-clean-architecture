package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate moq -out=../../mock/MockRepositoryOutbox.go -pkg=mock . OutboxRepository

type OutboxRepository interface {
	Insert(ctx context.Context, db *gorm.DB, outbox *entity.Outbox) error
	FindPending(ctx context.Context, db *gorm.DB, outboxes *entity.OutboxList, limit int) error
	MarkProduced(ctx context.Context, db *gorm.DB, ids []int64) error
}

var _ OutboxRepository = &OutboxRepositoryImpl{}

type OutboxRepositoryImpl struct {
	Cfg *config.Config
}

func NewOutboxRepository(cfg *config.Config) *OutboxRepositoryImpl {
	return &OutboxRepositoryImpl{
		Cfg: cfg,
	}
}

func (r *OutboxRepositoryImpl) Insert(ctx context.Context, db *gorm.DB, outbox *entity.Outbox) error {
	err := db.WithContext(ctx).Create(outbox).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*OutboxRepositoryImpl).Insert")
	}
	return nil
}

func (r *OutboxRepositoryImpl) FindPending(ctx context.Context, db *gorm.DB, outboxes *entity.OutboxList, limit int) error {
	err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where(column.Status.Eq(entity.OutboxStatusPending)).
		Order(column.CreatedAt.Str()).
		Limit(limit).
		Find(outboxes).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*OutboxRepositoryImpl).FindPending")
	}
	return nil
}

func (r *OutboxRepositoryImpl) MarkProduced(ctx context.Context, db *gorm.DB, ids []int64) error {
	err := db.WithContext(ctx).
		Model(&entity.Outbox{}).
		Where(column.ID.Eq(ids)).
		Update(column.Status.Str(), entity.OutboxStatusProduced).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*OutboxRepositoryImpl).MarkProduced")
	}
	return nil
}
