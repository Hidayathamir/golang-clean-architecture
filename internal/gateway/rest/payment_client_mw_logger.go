package rest

import (
	"context"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ PaymentClient = &PaymentClientMwLogger{}

type PaymentClientMwLogger struct {
	logger *logrus.Logger

	next PaymentClient
}

func NewPaymentClientMwLogger(logger *logrus.Logger, next PaymentClient) *PaymentClientMwLogger {
	return &PaymentClientMwLogger{
		logger: logger,
		next:   next,
	}
}

func (c *PaymentClientMwLogger) Refund(ctx context.Context, transactionID string) (bool, error) {
	ok, err := c.next.Refund(ctx, transactionID)

	fields := logrus.Fields{
		"transactionID": transactionID,
		"ok":            ok,
	}
	helper.Log(ctx, fields, err)

	return ok, err
}

func (c *PaymentClientMwLogger) GetStatus(ctx context.Context, transactionID string) (string, error) {
	status, err := c.next.GetStatus(ctx, transactionID)

	fields := logrus.Fields{
		"transactionID": transactionID,
		"status":        status,
	}
	helper.Log(ctx, fields, err)

	return status, err
}
