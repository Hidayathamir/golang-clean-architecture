package rest

import (
	"context"

	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/PaymentClient.go -pkg=mock . PaymentClient

type PaymentClient interface {
	Refund(ctx context.Context, transactionID string) (bool, error)
	GetStatus(ctx context.Context, transactionID string) (string, error)
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

func (c *PaymentClientImpl) Refund(ctx context.Context, transactionID string) (bool, error) {
	// TODO implement hit external rest api
	return true, nil
}

func (c *PaymentClientImpl) GetStatus(ctx context.Context, transactionID string) (string, error) {
	// TODO implement hit external rest api
	return "SUCCESS", nil
}
