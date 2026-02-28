package pkg

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitTracer initializes the OpenTelemetry TracerProvider with an OTLP gRPC exporter.
// It sets up the global TracerProvider and Propagator (W3C TraceContext) for the application.
//
// serviceName is the name of this microservice, env is the environment (staging/prod),
// and otlpEndpoint is the address of the OTLP collector (e.g., "jaeger:4317").
//
// It returns a shutdown function that should be executed via defer on the main entry point
// to ensure all spans are flushed to the exporter before the process terminates.
//
// Example:
//
//	shutdown := platform.InitTracer(ctx, "my-service", "prod", "otel-collector:4317")
//	defer shutdown(ctx)
func InitTracer(ctx context.Context, serviceName, env, otlpEndpoint string) func(context.Context) {
	// Create OTLP gRPC exporter to send traces to Jaeger/collector.
	conn, err := grpc.NewClient(otlpEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("Failed to create gRPC connection", "endpoint", otlpEndpoint, "error", err)
		os.Exit(1)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		slog.Error("Failed to create OTLP trace exporter", "error", err)
		os.Exit(1)
	}

	// Define the resource (metadata) for this service.
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		slog.Error("Failed to create resource", "error", err)
		os.Exit(1)
	}

	// Create the TracerProvider with the exporter and resource.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set the global TracerProvider so all otel instrumentation picks it up.
	otel.SetTracerProvider(tp)

	// Set the global Propagator to W3C TraceContext + Baggage.
	// This ensures that trace IDs are passed between services via HTTP headers.
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	slog.Info("OpenTelemetry initialized", "service", serviceName, "env", env, "endpoint", otlpEndpoint)

	// Return a shutdown function to flush remaining spans on app exit.
	return func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down tracer provider", "error", err)
		}
	}
}
