package telemetry

import (
	"context"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

var logger log.Logger = global.GetLoggerProvider().Logger(instrumentationScope)

func InitLogProvider(viperConfig *viper.Viper) (Stop, error) {
	exporter, err := newLogExporter(viperConfig)
	if err != nil {
		return nil, errkit.AddFuncName("telemetry.InitLogProvider", err)
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(viperConfig.GetString(configkey.AppName)),
			),
		),
	)

	global.SetLoggerProvider(lp)

	logger = lp.Logger(viperConfig.GetString(configkey.AppName))

	stop := func() {
		if err := lp.ForceFlush(context.Background()); err != nil {
			fmt.Println(errkit.AddFuncName("telemetry.InitLogProvider", err))
		}
		if err := lp.Shutdown(context.Background()); err != nil {
			fmt.Println(errkit.AddFuncName("telemetry.InitLogProvider", err))
		}
	}

	return stop, nil
}

func newLogExporter(viperConfig *viper.Viper) (sdklog.Exporter, error) {
	endpoint := viperConfig.GetString(configkey.TelemetryOTLPEndpoint)
	if endpoint == "" {
		err := fmt.Errorf("missing OTLP endpoint configuration")
		return nil, errkit.AddFuncName("telemetry.newLogExporter", err)
	}

	options := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(),
	}

	exporter, err := otlploggrpc.New(context.Background(), options...)
	if err != nil {
		return nil, errkit.AddFuncName("telemetry.newLogExporter", err)
	}

	return exporter, nil
}
