package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
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
	channelList, err := c.Next.GetChannelList(ctx)

	fields := logrus.Fields{
		"channelList": channelList,
	}
	logging.Log(ctx, fields, err)

	return channelList, err
}

func (c *SlackClientMwLogger) IsConnected(ctx context.Context) (bool, error) {
	isConnect, err := c.Next.IsConnected(ctx)

	fields := logrus.Fields{
		"isConnect": isConnect,
	}
	logging.Log(ctx, fields, err)

	return isConnect, err
}
