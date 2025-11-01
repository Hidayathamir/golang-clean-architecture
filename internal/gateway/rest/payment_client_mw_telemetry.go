package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ PaymentClient = &PaymentClientMwTelemetry{}

type PaymentClientMwTelemetry struct {
	Next PaymentClient
}

func NewPaymentClientMwTelemetry(next PaymentClient) *PaymentClientMwTelemetry {
	return &PaymentClientMwTelemetry{
		Next: next,
	}
}

func (c *PaymentClientMwTelemetry) Refund(ctx context.Context, transactionID string) (bool, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	ok, err := c.Next.Refund(ctx, transactionID)
	telemetry.RecordError(span, err)

	return ok, err
}

func (c *PaymentClientMwTelemetry) GetStatus(ctx context.Context, transactionID string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	status, err := c.Next.GetStatus(ctx, transactionID)
	telemetry.RecordError(span, err)

	return status, err
}
