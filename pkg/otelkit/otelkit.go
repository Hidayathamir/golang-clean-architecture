package otelkit

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

func ValidateAbleToExportSpan() {
	tracer := otel.Tracer("manual-validation-web")
	_, span := tracer.Start(context.Background(), "startup-check-web")
	span.SetAttributes(attribute.String("check", "success"))
	span.End()

	flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tp, ok := otel.GetTracerProvider().(*trace.TracerProvider); ok {
		err := tp.ForceFlush(flushCtx)
		if err != nil {
			err = errkit.SetMessage(err, "error export span, wait a little longer, or check is the collector ready")
			x.Logger.WithError(err).Panic()
		}
		x.Logger.Info("Successfully sent manual trace for web")
	}
}
