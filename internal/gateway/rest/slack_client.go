package rest

import (
	"context"

	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/SlackClient.go -pkg=mock . SlackClient

type SlackClient interface {
	GetChannelList(ctx context.Context) ([]string, error)
	IsConnected(ctx context.Context) (bool, error)
}

var _ SlackClient = &SlackClientImpl{}

type SlackClientImpl struct {
	Config *viper.Viper
}

func NewSlackClient(cfg *viper.Viper) *SlackClientImpl {
	return &SlackClientImpl{
		Config: cfg,
	}
}

func (c *SlackClientImpl) GetChannelList(ctx context.Context) ([]string, error) {
	// TODO implement hit external rest api
	return []string{"general", "random"}, nil
}

func (c *SlackClientImpl) IsConnected(ctx context.Context) (bool, error) {
	// TODO implement hit external rest api
	return true, nil
}
