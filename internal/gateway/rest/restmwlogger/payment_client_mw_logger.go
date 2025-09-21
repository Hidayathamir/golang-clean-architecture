package restmwlogger

import (
	"context"
	"golang-clean-architecture/internal/gateway/rest"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ rest.PaymentClient = &PaymentClientImpl{}

type PaymentClientImpl struct {
	logger *logrus.Logger

	next rest.PaymentClient
}

func NewPaymentClient(logger *logrus.Logger, next rest.PaymentClient) *PaymentClientImpl {
	return &PaymentClientImpl{
		logger: logger,
		next:   next,
	}
}

func (c *PaymentClientImpl) Refund(ctx context.Context, transactionID string) (bool, error) {
	ok, err := c.next.Refund(ctx, transactionID)

	fields := logrus.Fields{
		"transactionID": transactionID,
		"ok":            ok,
	}
	helper.Log(ctx, fields, err)

	return ok, err
}

func (c *PaymentClientImpl) GetStatus(ctx context.Context, transactionID string) (string, error) {
	status, err := c.next.GetStatus(ctx, transactionID)

	fields := logrus.Fields{
		"transactionID": transactionID,
		"status":        status,
	}
	helper.Log(ctx, fields, err)

	return status, err
}
