package httpComm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var _ (Api) = (*Client)(nil)

var (
	tracer            = otel.Tracer("github.com/oslokommune/httpcomm")
	traceCommonLabels = []attribute.KeyValue{
		attribute.String("language", "go"),
	}
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	headers    map[string]string
	tracing    bool
}

func NewClient(baseURL string, headers map[string]string, tracing bool) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	var httpClient *http.Client

	if tracing {
		httpClient = &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   30 * time.Second,
		}
	} else {
		httpClient = &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   30 * time.Second,
		}
	}

	return &Client{
		baseURL:    parsedURL,
		headers:    headers,
		httpClient: httpClient,
		tracing:    tracing,
	}, nil
}

func (c *Client) Post(ctx context.Context, resourceURL string, body io.Reader, contentType string, token *string) ([]byte, error) {
	var span trace.Span
	if c.tracing {
		ctx, span = tracer.Start(ctx, "post", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(traceCommonLabels...))
		defer span.End()
	}

	req, err := http.NewRequestWithContext(ctx, "POST", resourceURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	return c.callEndpoint(req, span)
}

func (c *Client) Put(ctx context.Context, resourceURL string, body io.Reader, contentType string, token *string) ([]byte, error) {
	var span trace.Span
	if c.tracing {
		ctx, span = tracer.Start(ctx, "put", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(traceCommonLabels...))
		defer span.End()
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", resourceURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	return c.callEndpoint(req, span)
}

func (c *Client) Delete(ctx context.Context, resourceURL string, token *string) ([]byte, error) {
	var span trace.Span
	if c.tracing {
		ctx, span = tracer.Start(ctx, "delete", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(traceCommonLabels...))
		defer span.End()
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", resourceURL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	return c.callEndpoint(req, span)
}

func (c *Client) Get(ctx context.Context, resourceURL string, token *string) ([]byte, error) {
	var span trace.Span
	if c.tracing {
		ctx, span = tracer.Start(ctx, "get", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(traceCommonLabels...))
		defer span.End()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", resourceURL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	return c.callEndpoint(req, span)
}

func (c *Client) callEndpoint(req *http.Request, span trace.Span) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	if c.tracing {
		span.SetAttributes(attribute.Int("http.status", statusCode))
		span.SetAttributes(attribute.String("http.method", "GET"))
		span.SetAttributes(attribute.String("http.url", req.URL.String()))
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("status code=%d, URL=%s", statusCode, req.URL.String())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
