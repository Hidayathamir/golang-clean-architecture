package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/SlackClient.go -pkg=mock . SlackClient

type SlackClient interface {
	GetChannelList(ctx context.Context, req model.SlackGetChannelListRequest) (model.SlackGetChannelListResponse, error)
	IsConnected(ctx context.Context, req model.SlackIsConnectedRequest) (model.SlackIsConnectedResponse, error)
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

func (c *SlackClientImpl) GetChannelList(ctx context.Context, req model.SlackGetChannelListRequest) (model.SlackGetChannelListResponse, error) {
	// TODO implement hit external rest api
	return model.SlackGetChannelListResponse{
		Channels: []string{"general", "random"},
	}, nil
}

func (c *SlackClientImpl) IsConnected(ctx context.Context, req model.SlackIsConnectedRequest) (model.SlackIsConnectedResponse, error) {
	// TODO implement hit external rest api
	return model.SlackIsConnectedResponse{
		Connected: true,
	}, nil
}
