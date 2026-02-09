package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/messaging/route"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	viperConfig := config.NewViper()
	x.SetupAll(viperConfig)
	db := config.NewDatabase(viperConfig)
	s3Client := config.NewS3Client(viperConfig)
	producer := config.NewKafkaProducer(viperConfig)

	usecases := dependency_injection.SetupUsecases(viperConfig, db, producer, s3Client)

	stopTraceProvider, err := telemetry.InitTraceProvider(viperConfig)
	panicIfErr(err)
	defer stopTraceProvider()

	validateAbleToExportSpan()

	stopLogProvider, err := telemetry.InitLogProvider(viperConfig)
	panicIfErr(err)
	defer stopLogProvider()

	ctx, cancel := context.WithCancel(context.Background())

	x.Logger.Info("Starting worker service")
	go route.ConsumeUserFollowedEventForNotification(ctx, viperConfig, usecases)
	go route.ConsumeUserFollowedEventForUpdateCount(ctx, viperConfig, usecases)
	go route.SetupImageUploadedConsumer(ctx, viperConfig, usecases)
	go route.ConsumeImageLikedEventForNotification(ctx, viperConfig, usecases)
	go route.ConsumeImageLikedEventForUpdateCount(ctx, viperConfig, usecases)
	go route.ConsumeImageCommentedEventForNotification(ctx, viperConfig, usecases)
	go route.ConsumeImageCommentedEventForUpdateCount(ctx, viperConfig, usecases)
	go route.SetupNotifConsumer(ctx, viperConfig, usecases)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	s := <-terminateSignals
	x.Logger.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)

	x.Logger.Info("canceling")
	cancel()
	x.Logger.Info("canceled")

	x.Logger.Info("wait for all consumer to finish processing, 5 second")
	time.Sleep(5 * time.Second)
	x.Logger.Info("done waiting")

	x.Logger.Info("end process of worker")
}

func validateAbleToExportSpan() {
	tracer := otel.Tracer("manual-validation-worker")
	_, span := tracer.Start(context.Background(), "startup-check-worker")
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
		x.Logger.Info("Successfully sent manual trace for worker")
	}
}

func panicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
