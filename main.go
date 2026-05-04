package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func main() {
	logger, cleanup := initOtelLogger()
	defer cleanup()

	runApp(logger)
}

// initOtelLogger creates an OpenTelemetry-enabled logger
// sends logs directly to the OTel collector via gRPC
func initOtelLogger() (*slog.Logger, func()) {
	fmt.Println("Running in OTLP mode")
	fmt.Println("   - Logs sent to OpenTelemetry Collector (localhost:4317)")
	fmt.Println("   - Full OpenTelemetry SDK integration")
	fmt.Println("   - Collector batches and forwards to Splunk")
	fmt.Println()

	ctx := context.Background()

	// create OTLP log exporter
	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint("localhost:4317"),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Failed to create OTLP exporter: %v\n   Make sure OTel collector is running!", err)
	}

	// create resource (service metadata)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("my-app"),
			semconv.ServiceVersion("1.0.0"),
			semconv.DeploymentEnvironment("development"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// create logger provider with batch processor
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)

	// set as global provider
	global.SetLoggerProvider(loggerProvider)

	// create slog handler that bridges to OTel
	otelHandler := otelslog.NewHandler("my-app", otelslog.WithLoggerProvider(loggerProvider))

	// create the slog logger
	logger := slog.New(otelHandler)

	// cleanup function to flush logs
	cleanup := func() {
		fmt.Println("\n Flushing logs to collector...")
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down logger provider: %v", err)
		}
		time.Sleep(2 * time.Second)
		fmt.Println("Logs sent to OpenTelemetry Collector!")
		fmt.Println("Check Splunk at http://localhost:8000")
		fmt.Println("Search: index=* source=\"my-app\"")
	}

	return logger, cleanup
}

func runApp(logger *slog.Logger) {
	logger.Info("application started",
		slog.String("version", "1.0.0"),
	)

	logger.Info("user login",
		slog.String("user_id", "user_123"),
		slog.String("username", "john"),
		slog.String("action", "login_success"),
	)

	logger.Error("failed to connect to database",
		slog.String("error", "connection timeout"),
		slog.String("database", "postgres"),
		slog.Int("retry_count", 3),
	)

	logger.Warn("high memory usage detected",
		slog.Float64("memory_mb", 1024.5),
		slog.Float64("threshold_mb", 1000.0),
	)

	logger.Info("application stopped")
}
