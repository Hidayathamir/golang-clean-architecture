package entity

import (
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"gorm.io/gorm"
)

type Outbox struct {
	ID           int64          `gorm:"column:id;primaryKey"`
	Topic        string         `gorm:"column:topic"`
	Key          string         `gorm:"column:key"`
	Payload      []byte         `gorm:"column:payload"`
	TraceContext string         `gorm:"column:trace_context"`
	Status       string         `gorm:"column:status"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (o *Outbox) TableName() string {
	return table.Outbox
}

type OutboxList []Outbox

const (
	OutboxStatusPending  = "pending"
	OutboxStatusProduced = "produced"
)
