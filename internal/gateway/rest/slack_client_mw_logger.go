package rest

import (
	"context"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ SlackClient = &SlackClientMwLogger{}

type SlackClientMwLogger struct {
	logger *logrus.Logger

	next SlackClient
}

func NewSlackClientMwLogger(logger *logrus.Logger, next SlackClient) *SlackClientMwLogger {
	return &SlackClientMwLogger{
		logger: logger,
		next:   next,
	}
}

func (c *SlackClientMwLogger) GetChannelList(ctx context.Context) ([]string, error) {
	channelList, err := c.next.GetChannelList(ctx)

	fields := logrus.Fields{
		"channelList": channelList,
	}
	helper.Log(ctx, fields, err)

	return channelList, err
}

func (c *SlackClientMwLogger) IsConnected(ctx context.Context) (bool, error) {
	isConnect, err := c.next.IsConnected(ctx)

	fields := logrus.Fields{
		"isConnect": isConnect,
	}
	helper.Log(ctx, fields, err)

	return isConnect, err
}
