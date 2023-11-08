package httpcomm

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel/trace"
)

func Decode[T any](ctx context.Context, responseBody []byte, tracing bool) (*T, error) {
	if tracing {
		_, span := tracer.Start(ctx, "decode-json", trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(traceCommonLabels...))
		defer span.End()
	}

	var message T
	err := json.Unmarshal([]byte(responseBody), &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
