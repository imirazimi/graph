package tracing

import (
    "context"
    "log"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer() func(context.Context) error {
    exporter, err := stdouttrace.New(
        stdouttrace.WithPrettyPrint(),
    )
    if err != nil {
        log.Fatalf("failed to initialize stdouttrace exporter: %v", err)
    }

    provider := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
    )

    otel.SetTracerProvider(provider)

    return provider.Shutdown
}