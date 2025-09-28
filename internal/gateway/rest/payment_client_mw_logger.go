package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/helper"
	"github.com/sirupsen/logrus"
)

var _ PaymentClient = &PaymentClientMwLogger{}

type PaymentClientMwLogger struct {
	Next PaymentClient
}

func NewPaymentClientMwLogger(next PaymentClient) *PaymentClientMwLogger {
	return &PaymentClientMwLogger{
		Next: next,
	}
}

func (c *PaymentClientMwLogger) Refund(ctx context.Context, transactionID string) (bool, error) {
	ok, err := c.Next.Refund(ctx, transactionID)

	fields := logrus.Fields{
		"transactionID": transactionID,
		"ok":            ok,
	}
	helper.Log(ctx, fields, err)

	return ok, err
}

func (c *PaymentClientMwLogger) GetStatus(ctx context.Context, transactionID string) (string, error) {
	status, err := c.Next.GetStatus(ctx, transactionID)

	fields := logrus.Fields{
		"transactionID": transactionID,
		"status":        status,
	}
	helper.Log(ctx, fields, err)

	return status, err
}
