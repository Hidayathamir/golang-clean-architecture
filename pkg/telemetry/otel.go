package telemetry

import (
	"context"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/caller"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

const defaultTracerName = "github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"

var tracer trace.Tracer = otel.Tracer(defaultTracerName)

type Stop func()

func Init(viperConfig *viper.Viper) (Stop, error) {
	exporter, err := newOTLPExporter(context.Background(), viperConfig)
	if err != nil {
		return nil, errkit.AddFuncName("telemetry.Init", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(viperConfig.GetString(configkey.AppName)),
			)),
	)

	otel.SetTracerProvider(tp)

	tracer = tp.Tracer(viperConfig.GetString(configkey.AppName))

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	stopTelemetry := func() {
		if err := tp.ForceFlush(context.Background()); err != nil {
			fmt.Println(errkit.AddFuncName("telemetry.Init", err))
		}
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Println(errkit.AddFuncName("telemetry.Init", err))
		}
	}

	return stopTelemetry, nil
}

func newOTLPExporter(ctx context.Context, viperConfig *viper.Viper) (sdktrace.SpanExporter, error) {
	endpoint := viperConfig.GetString(configkey.TelemetryOTLPEndpoint)
	if endpoint == "" {
		err := fmt.Errorf("missing OTLP endpoint configuration")
		return nil, errkit.AddFuncName("telemetry.newOTLPExporter", err)
	}

	clientOptions := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	}

	exporter, err := otlptracegrpc.New(ctx, clientOptions...)
	if err != nil {
		return nil, errkit.AddFuncName("telemetry.newOTLPExporter", err)
	}

	return exporter, nil
}

func StartController(ctx *fiber.Ctx) trace.Span {
	newctx, span := tracer.Start(ctx.UserContext(), caller.FuncName(caller.WithSkip(1)))
	ctx.SetUserContext(newctx)
	return span
}

func StartConsumer(message *sarama.ConsumerMessage) (context.Context, trace.Span) {
	ctx := ExtractCtxFromConsumerMessage(message)

	return tracer.Start(ctx, caller.FuncName(caller.WithSkip(1)))
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
