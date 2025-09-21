package rest

import "context"

//go:generate moq -out=../../mock/SlackClient.go -pkg=mock . SlackClient

type SlackClient interface {
	GetChannelList(ctx context.Context) ([]string, error)
	IsConnected(ctx context.Context) (bool, error)
}

var _ SlackClient = &SlackClientImpl{}

type SlackClientImpl struct {
}

func NewSlackClient() *SlackClientImpl {
	return &SlackClientImpl{}
}

func (c *SlackClientImpl) GetChannelList(ctx context.Context) ([]string, error) {
	// TODO implement hit external rest api
	return []string{"general", "random"}, nil
}

func (c *SlackClientImpl) IsConnected(ctx context.Context) (bool, error) {
	// TODO implement hit external rest api
	return true, nil
}
