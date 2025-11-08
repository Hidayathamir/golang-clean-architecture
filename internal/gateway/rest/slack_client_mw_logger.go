package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ SlackClient = &SlackClientMwLogger{}

type SlackClientMwLogger struct {
	Next SlackClient
}

func NewSlackClientMwLogger(next SlackClient) *SlackClientMwLogger {
	return &SlackClientMwLogger{
		Next: next,
	}
}

func (c *SlackClientMwLogger) GetChannelList(ctx context.Context, req model.SlackGetChannelListRequest) (model.SlackGetChannelListResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := c.Next.GetChannelList(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"channels": res.Channels,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}

func (c *SlackClientMwLogger) IsConnected(ctx context.Context, req model.SlackIsConnectedRequest) (model.SlackIsConnectedResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := c.Next.IsConnected(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"isConnect": res.Connected,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}
