package telemetry

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/caller"
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationScope = "github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"

type Stop func()

func StartController(ctx *fiber.Ctx) trace.Span {
	newctx, span := tracer.Start(ctx.UserContext(), caller.FuncName(caller.WithSkip(1)))
	ctx.SetUserContext(newctx)
	return span
}

func StartConsumer(message *sarama.ConsumerMessage) (context.Context, trace.Span) {
	ctx := ExtractCtxFromConsumerMessage(message)

	return tracer.Start(ctx, caller.FuncName(caller.WithSkip(1)))
}

func StartNew() (context.Context, trace.Span) {
	return tracer.Start(context.Background(), caller.FuncName(caller.WithSkip(1)))
}

func ExtractCtxFromConsumerMessage(message *sarama.ConsumerMessage) context.Context {
	carrier := otelsarama.NewConsumerMessageCarrier(message)

	ctx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

	return ctx
}

func InjectCtxToProducerMessage(ctx context.Context, message *sarama.ProducerMessage) {
	carrier := otelsarama.NewProducerMessageCarrier(message)

	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

func Start(ctx context.Context) (context.Context, trace.Span) {
	return tracer.Start(ctx, caller.FuncName(caller.WithSkip(1)))
}

func RecordError(span trace.Span, err error) {
	if err == nil {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	sc := span.SpanContext()
	if !sc.IsValid() || !sc.TraceID().IsValid() {
		return ""
	}

	return sc.TraceID().String()
}

func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	sc := span.SpanContext()
	if !sc.IsValid() || !sc.SpanID().IsValid() {
		return ""
	}

	return sc.SpanID().String()
}
