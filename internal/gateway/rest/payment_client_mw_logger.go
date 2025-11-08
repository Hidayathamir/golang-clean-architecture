package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
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
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	ok, err := c.Next.Refund(ctx, transactionID)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"transactionID": transactionID,
		"ok":            ok,
	}
	x.LogMw(ctx, fields, err)

	return ok, err
}

func (c *PaymentClientMwLogger) GetStatus(ctx context.Context, transactionID string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	status, err := c.Next.GetStatus(ctx, transactionID)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"transactionID": transactionID,
		"status":        status,
	}
	x.LogMw(ctx, fields, err)

	return status, err
}
