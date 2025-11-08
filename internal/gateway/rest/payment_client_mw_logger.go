package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
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

func (c *PaymentClientMwLogger) Refund(ctx context.Context, req model.PaymentRefundRequest) (model.PaymentRefundResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := c.Next.Refund(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"transactionID": req.TransactionID,
		"success":       res.Success,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}

func (c *PaymentClientMwLogger) GetStatus(ctx context.Context, req model.PaymentGetStatusRequest) (model.PaymentGetStatusResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := c.Next.GetStatus(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"transactionID": req.TransactionID,
		"status":        res.Status,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}
