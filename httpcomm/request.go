package httpcomm

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer            = otel.Tracer("github.com/oslokommune/httpcomm")
	traceCommonLabels = []attribute.KeyValue{
		attribute.String("language", "go"),
	}
)

// Interface for performing HTTP calls, e.g. a `http.Client`
type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

func createRequest(ctx context.Context, httpRequest HTTPRequest) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, httpRequest.Method, httpRequest.Url, httpRequest.Body)
	if err != nil {
		return nil, err
	}

	for key, value := range httpRequest.Headers {
		req.Header.Set(key, value)
	}

	if httpRequest.Token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *httpRequest.Token))
	}

	return req, nil
}

func Call(ctx context.Context, httpClient HttpDoer, httpRequest HTTPRequest) (*HTTPResponse, error) {
	var span trace.Span
	if httpRequest.Tracing {
		_, span = tracer.Start(ctx, httpRequest.Method, trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(traceCommonLabels...))
		defer span.End()
	}

	req, err := createRequest(ctx, httpRequest)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	if httpRequest.Tracing {
		span.SetAttributes(attribute.Int("http.status", statusCode))
		span.SetAttributes(attribute.String("http.method", "GET"))
		span.SetAttributes(attribute.String("http.url", req.URL.String()))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, &HTTPError{
			Body:       string(body),
			StatusCode: statusCode,
		}
	}

	return &HTTPResponse{
		StatusCode: statusCode,
		Body:       string(body),
	}, nil
}
