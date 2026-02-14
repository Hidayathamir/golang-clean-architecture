package telemetry

import (
	"context"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

var logger log.Logger = global.GetLoggerProvider().Logger(instrumentationScope)

func InitLogProvider(cfg *config.Config) Stop {
	exporter := newLogExporter(cfg)

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.GetAppName()),
			),
		),
	)

	global.SetLoggerProvider(lp)

	logger = lp.Logger(cfg.GetAppName())

	stop := func() {
		err := lp.ForceFlush(context.Background())
		if err != nil {
			fmt.Println(errkit.AddFuncName(err))
		}
		err = lp.Shutdown(context.Background())
		if err != nil {
			fmt.Println(errkit.AddFuncName(err))
		}
	}

	return stop
}

func newLogExporter(cfg *config.Config) sdklog.Exporter {
	endpoint := cfg.GetTelemetryOTLPEndpoint()
	if endpoint == "" {
		err := fmt.Errorf("missing OTLP endpoint configuration")
		panic(err)
	}

	options := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(),
	}

	exporter, err := otlploggrpc.New(context.Background(), options...)
	if err != nil {
		panic(err)
	}

	return exporter
}
