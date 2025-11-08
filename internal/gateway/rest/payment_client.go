package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/PaymentClient.go -pkg=mock . PaymentClient

type PaymentClient interface {
	Refund(ctx context.Context, req model.PaymentRefundRequest) (model.PaymentRefundResponse, error)
	GetStatus(ctx context.Context, req model.PaymentGetStatusRequest) (model.PaymentGetStatusResponse, error)
}

var _ PaymentClient = &PaymentClientImpl{}

type PaymentClientImpl struct {
	Config *viper.Viper
}

func NewPaymentClient(cfg *viper.Viper) *PaymentClientImpl {
	return &PaymentClientImpl{
		Config: cfg,
	}
}

func (c *PaymentClientImpl) Refund(ctx context.Context, req model.PaymentRefundRequest) (model.PaymentRefundResponse, error) {
	// TODO implement hit external rest api
	return model.PaymentRefundResponse{
		Success: true,
	}, nil
}

func (c *PaymentClientImpl) GetStatus(ctx context.Context, req model.PaymentGetStatusRequest) (model.PaymentGetStatusResponse, error) {
	// TODO implement hit external rest api
	return model.PaymentGetStatusResponse{
		Status: "SUCCESS",
	}, nil
}
