package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ ChatRepository = &ChatRepositoryMwLogger{}

type ChatRepositoryMwLogger struct {
	Next ChatRepository
}

func NewChatRepositoryMwLogger(next ChatRepository) *ChatRepositoryMwLogger {
	return &ChatRepositoryMwLogger{
		Next: next,
	}
}

func (r *ChatRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, chat *entity.Chat) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.Create(ctx, db, chat)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"chat": chat,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *ChatRepositoryMwLogger) FindThread(ctx context.Context, db *gorm.DB, chatList *entity.ChatList, userID, peerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindThread(ctx, db, chatList, userID, peerID, cursorCreatedAt, cursorID, limit)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"user_id": userID,
		"peer_id": peerID,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *ChatRepositoryMwLogger) ListConversations(ctx context.Context, db *gorm.DB, conversations *entity.ConversationList, userID int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.ListConversations(ctx, db, conversations, userID)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"user_id": userID,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *ChatRepositoryMwLogger) MarkThreadRead(ctx context.Context, db *gorm.DB, userID, peerID int64, readAt time.Time) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	rows, err := retrykit.DBRetryWithData(ctx, func() (int64, error) {
		return r.Next.MarkThreadRead(ctx, db, userID, peerID, readAt)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"user_id": userID,
		"peer_id": peerID,
	}
	logkit.LogMw(ctx, fields, err)

	return rows, err
}

func (r *ChatRepositoryMwLogger) CountUnread(ctx context.Context, db *gorm.DB, userID int64) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	total, err := retrykit.DBRetryWithData(ctx, func() (int64, error) {
		return r.Next.CountUnread(ctx, db, userID)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"user_id": userID,
	}
	logkit.LogMw(ctx, fields, err)

	return total, err
}
