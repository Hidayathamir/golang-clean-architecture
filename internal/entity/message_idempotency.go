package entity

import (
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
)

type MessageIdempotency struct {
	IdempotencyKey string    `gorm:"column:idempotency_key;primaryKey"`
	Topic          string    `gorm:"column:topic"`
	Partition      int32     `gorm:"column:partition"`
	RecordOffset   int64     `gorm:"column:record_offset"`
	ProcessedAt    time.Time `gorm:"column:processed_at;autoCreateTime"`
}

func (m *MessageIdempotency) TableName() string {
	return table.MessageIdempotency
}
