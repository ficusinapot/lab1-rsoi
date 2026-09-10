package tracing

import (
	"context"

	"github.com/ficusinapot/ds/internal/config"

	"github.com/joomcode/errorx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewTracerProvider(ctx context.Context, appName string, appVersion string, cfg config.OpenTelemetryConfig) (*sdktrace.TracerProvider, error) {
	if cfg.Address == "" {
		provider := sdktrace.NewTracerProvider()
		otel.SetTracerProvider(provider)
		setPropagator()

		return provider, nil
	}

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.Address),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, errorx.ExternalError.Wrap(err, "create OTLP trace exporter")
	}

	serviceName := cfg.ServiceName
	if serviceName == "" {
		serviceName = appName
	}

	batcherOptions := make([]sdktrace.BatchSpanProcessorOption, 0, 1)
	if cfg.QueueSize > 0 {
		batcherOptions = append(batcherOptions, sdktrace.WithMaxQueueSize(cfg.QueueSize))
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter, batcherOptions...),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(appVersion),
		)),
	)

	otel.SetTracerProvider(provider)
	setPropagator()

	return provider, nil
}

func setPropagator() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}
