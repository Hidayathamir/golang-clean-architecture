package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ SlackClient = &SlackClientMwTelemetry{}

type SlackClientMwTelemetry struct {
	Next SlackClient
}

func NewSlackClientMwTelemetry(next SlackClient) *SlackClientMwTelemetry {
	return &SlackClientMwTelemetry{
		Next: next,
	}
}

func (c *SlackClientMwTelemetry) GetChannelList(ctx context.Context) ([]string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	list, err := c.Next.GetChannelList(ctx)
	telemetry.RecordError(span, err)

	return list, err
}

func (c *SlackClientMwTelemetry) IsConnected(ctx context.Context) (bool, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	ok, err := c.Next.IsConnected(ctx)
	telemetry.RecordError(span, err)

	return ok, err
}
