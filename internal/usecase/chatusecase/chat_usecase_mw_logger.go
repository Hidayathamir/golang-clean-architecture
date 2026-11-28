package chatusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ ChatUsecase = &ChatUsecaseMwLogger{}

type ChatUsecaseMwLogger struct {
	Next ChatUsecase
}

func NewChatUsecaseMwLogger(next ChatUsecase) *ChatUsecaseMwLogger {
	return &ChatUsecaseMwLogger{
		Next: next,
	}
}

func (u *ChatUsecaseMwLogger) Send(ctx context.Context, req dto.SendChatRequest) (dto.ChatResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Send(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	logkit.LogMw(ctx, fields, err)

	return res, err
}

func (u *ChatUsecaseMwLogger) ListConversations(ctx context.Context) (dto.ConversationsResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.ListConversations(ctx)
	telemetry.RecordError(span, err)

	logkit.LogMw(ctx, logrus.Fields{}, err)

	return res, err
}

func (u *ChatUsecaseMwLogger) GetThread(ctx context.Context, req dto.GetThreadRequest) (dto.ThreadResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.GetThread(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	logkit.LogMw(ctx, fields, err)

	return res, err
}

func (u *ChatUsecaseMwLogger) MarkRead(ctx context.Context, req dto.MarkChatReadRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.MarkRead(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (u *ChatUsecaseMwLogger) CountUnread(ctx context.Context) (dto.UnreadChatCountResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.CountUnread(ctx)
	telemetry.RecordError(span, err)

	logkit.LogMw(ctx, logrus.Fields{}, err)

	return res, err
}
