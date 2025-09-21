package restmwlogger

import (
	"context"
	"golang-clean-architecture/internal/gateway/rest"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ rest.SlackClient = &SlackClientImpl{}

type SlackClientImpl struct {
	logger *logrus.Logger

	next rest.SlackClient
}

func NewSlackClient(logger *logrus.Logger, next rest.SlackClient) *SlackClientImpl {
	return &SlackClientImpl{
		logger: logger,
		next:   next,
	}
}

func (c *SlackClientImpl) GetChannelList(ctx context.Context) ([]string, error) {
	channelList, err := c.next.GetChannelList(ctx)

	fields := logrus.Fields{
		"channelList": channelList,
	}
	helper.Log(ctx, fields, err)

	return channelList, err
}

func (c *SlackClientImpl) IsConnected(ctx context.Context) (bool, error) {
	isConnect, err := c.next.IsConnected(ctx)

	fields := logrus.Fields{
		"isConnect": isConnect,
	}
	helper.Log(ctx, fields, err)

	return isConnect, err
}
