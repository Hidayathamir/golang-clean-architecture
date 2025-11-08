package rest

import (
	"context"

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

func (c *SlackClientMwLogger) GetChannelList(ctx context.Context) ([]string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	channelList, err := c.Next.GetChannelList(ctx)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"channelList": channelList,
	}
	x.LogMw(ctx, fields, err)

	return channelList, err
}

func (c *SlackClientMwLogger) IsConnected(ctx context.Context) (bool, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	isConnect, err := c.Next.IsConnected(ctx)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"isConnect": isConnect,
	}
	x.LogMw(ctx, fields, err)

	return isConnect, err
}
