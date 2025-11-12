package telemetry

import (
	"context"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer = otel.Tracer(instrumentationScope)

func InitTraceProvider(viperConfig *viper.Viper) (Stop, error) {
	exporter, err := newSpanExporter(viperConfig)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(viperConfig.GetString(configkey.AppName)),
			),
		),
	)

	otel.SetTracerProvider(tp)

	tracer = tp.Tracer(viperConfig.GetString(configkey.AppName))

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	stop := func() {
		if err := tp.ForceFlush(context.Background()); err != nil {
			fmt.Println(errkit.AddFuncName(err))
		}
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Println(errkit.AddFuncName(err))
		}
	}

	return stop, nil
}

func newSpanExporter(viperConfig *viper.Viper) (sdktrace.SpanExporter, error) {
	endpoint := viperConfig.GetString(configkey.TelemetryOTLPEndpoint)
	if endpoint == "" {
		err := fmt.Errorf("missing OTLP endpoint configuration")
		return nil, errkit.AddFuncName(err)
	}

	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	}

	exporter, err := otlptracegrpc.New(context.Background(), options...)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return exporter, nil
}
