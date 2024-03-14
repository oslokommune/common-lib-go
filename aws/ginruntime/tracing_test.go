package ginruntime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Attributes []attribute.KeyValue

func (this Attributes) Get(key attribute.Key) attribute.Value {
	for _, kv := range this {
		if kv.Key == key {
			return kv.Value
		}
	}
	return attribute.Value{}
}
func (this Attributes) Has(key attribute.Key) bool {
	return this.Get(key) != attribute.Value{}
}

type InterceptingTraceExporter struct {
	Exporter trace.SpanExporter
	Callback func([]trace.ReadOnlySpan)
}

func (this *InterceptingTraceExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	this.Callback(spans)
	return this.Exporter.ExportSpans(ctx, spans)
}

func (this *InterceptingTraceExporter) Shutdown(ctx context.Context) error {
	return this.Exporter.Shutdown(ctx)
}

func NewInterceptingTracerProvider(callback func([]trace.ReadOnlySpan)) *trace.TracerProvider {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create stdouttrace exporter")
	}
	return trace.NewTracerProvider(
		trace.WithBatcher(&InterceptingTraceExporter{exporter, callback}),
		trace.WithSampler(trace.AlwaysSample()),
	)
}

func TestTracingEnabled_ReturnsFalse_ByDefault(t *testing.T) {
	engine := New(context.Background())
	assert.False(t, engine.TracingEnabled())
}

func TestTracingEnabled_ReturnsTrue_AfterEnableCustomTracingCalled(t *testing.T) {
	interceptor := NewInterceptingTracerProvider(func(spans []trace.ReadOnlySpan) {})
	engine := New(context.Background(), WithTracing("test", interceptor, &xray.Propagator{}))
	assert.True(t, engine.TracingEnabled())
}

func TestMiddleware_TracesHttpServerSpans(t *testing.T) {
	ctx := context.Background()
	spans := []trace.ReadOnlySpan{}
	interceptor := NewInterceptingTracerProvider(func(exportedSpans []trace.ReadOnlySpan) {
		spans = append(spans, exportedSpans...)
	})

	engine := New(context.Background(), WithTracing("test", interceptor, &xray.Propagator{}))
	engine.AddRoute(nil, "/foo/:bar", GET, nil, func(c *gin.Context) {
		c.JSON(200, "bar")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo/test", nil)
	engine.ServerHttp(res, req)
	if err := interceptor.ForceFlush(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to flush spans")
	}
	assert.Equal(t, 1, len(spans))
}

func TestMiddleware_TracesHttpServerSpansWithExpectedAttributes(t *testing.T) {
	ctx := context.Background()
	service := "test"
	spans := []trace.ReadOnlySpan{}
	interceptor := NewInterceptingTracerProvider(func(exportedSpans []trace.ReadOnlySpan) {
		spans = append(spans, exportedSpans...)
	})

	engine := New(context.Background(), WithTracing(service, interceptor, &xray.Propagator{}))

	engine.AddRoute(nil, "/foo/:bar", GET, nil, func(c *gin.Context) {
		c.JSON(200, "bar")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo/test?q=true", nil)
	engine.ServerHttp(res, req)
	if err := interceptor.ForceFlush(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to flush spans")
	}

	assert.Equal(t, 1, len(spans))
	assert.Equal(t, service, spans[0].Name())
	assert.Equal(t, oteltrace.SpanKindServer, spans[0].SpanKind())
	attributes := Attributes(spans[0].Attributes())
	assert.Equal(t, "GET", attributes.Get(semconv.HTTPMethodKey).AsString())
	assert.Equal(t, int64(200), attributes.Get(semconv.HTTPStatusCodeKey).AsInt64())
	// http.target is deprecated and should be replaced with url.path and url.query
	if attributes.Has(semconv.HTTPTargetKey) {
		// otelgin doesn't include the query in http.target
		assert.Equal(t, "/foo/test", attributes.Get(semconv.HTTPTargetKey).AsString())
	} else {
		assert.Equal(t, "/foo/test", attributes.Get(semconv.URLPathKey).AsString())
		assert.Equal(t, "q=true", attributes.Get(semconv.URLQueryKey).AsString())
	}

	// otelgin unfortunately uses the span name as the http.route value.
	// so by annotating the span with the service name we lose out on that information.
	// assert.Equal(t, "/foo/:bar", attributes.Get(semconv.HTTPRouteKey).AsString())
}

func TestEnableTracing_TracesLambdaInvocationWithExpectedAttributes(t *testing.T) {
	requestId := "requestId"
	service := "test"
	invocationArn := "invocationArn"
	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test")

	ctx := lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{
		AwsRequestID:       requestId,
		InvokedFunctionArn: invocationArn,
		ClientContext:      lambdacontext.ClientContext{},
		Identity:           lambdacontext.CognitoIdentity{}},
	)

	spans := []trace.ReadOnlySpan{}
	interceptor := NewInterceptingTracerProvider(func(exportedSpans []trace.ReadOnlySpan) {
		spans = append(spans, exportedSpans...)
	})

	engine := New(ctx, WithTracing(service, interceptor, &xray.Propagator{}))
	engine.AddRoute(nil, "/foo/:bar", GET, nil, func(c *gin.Context) {
		c.JSON(200, "bar")
	})

	proxy := engine.lambdaProxy()
	req := events.APIGatewayV2HTTPRequest{
		RawPath: "/foo/test",
	}
	_, err := proxy(ctx, req)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to invoke lambda proxy")
	}

	if err := interceptor.ForceFlush(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to flush spans")
	}

	// This includes the otelgin HTTP server span
	assert.Equal(t, 2, len(spans))

	attributes := Attributes(spans[0].Attributes())

	if !attributes.Has(semconv.FaaSInvocationIDKey) {
		attributes = Attributes(spans[1].Attributes())
	}

	assert.Equal(t, service, spans[0].Name())
	assert.Equal(t, requestId, attributes.Get(semconv.FaaSInvocationIDKey).AsString())
	assert.Equal(t, invocationArn, attributes.Get(semconv.AWSLambdaInvokedARNKey).AsString())
}

func TestEnableTracing_TracesOtelHttpClientSpansAsSubsegments(t *testing.T) {

	ctx := context.Background()
	service := "test"
	spans := []trace.ReadOnlySpan{}
	interceptor := NewInterceptingTracerProvider(func(exportedSpans []trace.ReadOnlySpan) {
		spans = append(spans, exportedSpans...)
	})

	engine := New(ctx, WithTracing(service, interceptor, &xray.Propagator{}))
	engine.AddRoute(nil, "/foo/:bar", GET, nil, func(c *gin.Context) {
		c.JSON(200, "bar")
		_, _ = otelhttp.Get(c.Request.Context(), "https://test")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo/test?", nil)
	engine.ServerHttp(res, req)
	if err := interceptor.ForceFlush(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to flush spans")
	}

	// This includes the otelgin HTTP server span
	assert.Equal(t, 2, len(spans))

	httpClientSpan := spans[0]
	httpServerSpan := spans[1]
	if httpClientSpan.SpanKind() != oteltrace.SpanKindClient {
		httpClientSpan = spans[1]
		httpServerSpan = spans[0]
	}
	attributes := Attributes(httpClientSpan.Attributes())
	assert.Equal(t, "GET", attributes.Get(semconv.HTTPMethodKey).AsString())
	assert.Equal(t, "https://test", attributes.Get(semconv.HTTPURLKey).AsString())
	assert.Equal(t, httpServerSpan.SpanContext().SpanID(), httpClientSpan.Parent().SpanID())
}

func TestEnableTracing_TracesHttpClientSpansAsSegments_WhenUsingOtelTransport(t *testing.T) {

	ctx := context.Background()
	service := "test"
	spans := []trace.ReadOnlySpan{}
	interceptor := NewInterceptingTracerProvider(func(exportedSpans []trace.ReadOnlySpan) {
		spans = append(spans, exportedSpans...)
	})

	engine := New(ctx, WithTracing(service, interceptor, &xray.Propagator{}))

	engine.AddRoute(nil, "/foo/:bar", GET, nil, func(c *gin.Context) {
		c.JSON(200, "bar")

		client := &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   time.Microsecond,
		}

		_, _ = client.Get("https://test")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo/test?", nil)
	engine.ServerHttp(res, req)
	if err := interceptor.ForceFlush(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to flush spans")
	}

	// This includes the otelgin HTTP server span
	assert.Equal(t, 2, len(spans))

	span := spans[0]
	if span.SpanKind() != oteltrace.SpanKindClient {
		span = spans[1]
	}

	attributes := Attributes(span.Attributes())
	assert.Equal(t, "GET", attributes.Get(semconv.HTTPMethodKey).AsString())
	assert.Equal(t, "https://test", attributes.Get(semconv.HTTPURLKey).AsString())
}
