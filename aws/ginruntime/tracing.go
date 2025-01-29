package ginruntime

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func (e *GinEngine) oTelLambdaOptions() []otellambda.Option {
	if !e.TracingEnabled() {
		log.Warn().Msg("tracing is not enabled - no otellambda options are available")
		return []otellambda.Option{}
	}

	return []otellambda.Option{
		xrayconfig.WithEventToCarrier(),
		otellambda.WithPropagator(e.propagator),
		otellambda.WithFlusher(e.tp),
		otellambda.WithTracerProvider(e.tp),
	}
}

func (e *GinEngine) otelGinOptions(service string) []otelgin.Option {
	if !e.TracingEnabled() {
		log.Warn().Msg("tracing is not enabled - no otelgin options are available")
		return []otelgin.Option{}
	}

	return []otelgin.Option{
		otelgin.WithSpanNameFormatter(func(h *http.Request) string { return service }),
		otelgin.WithTracerProvider(e.tp),
		otelgin.WithPropagators(e.propagator),
	}
}

func (e *GinEngine) configureOpenTelemetry() {
	otel.SetTracerProvider(e.tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			e.propagator,
		),
	)
}

func (e *GinEngine) useHttpServerSpanMiddleware(service string) {
	if service == "" {
		log.Fatal().Msg("missing service name to use for tracing")
	}

	e.Use(otelgin.Middleware(service, e.otelGinOptions(service)...))
}

// Extracts trace ID from the logging context and writes it to the `X-Amzn-Trace-Id` key in the log event.
//
// Usage:
//
//	log.Logger = log.Logger.Hook(ginruntime.XAmznTraceIdLoggerHook{})
type XAmznTraceIdLoggerHook struct{}

const xAmznTraceIdHeader = "X-Amzn-Trace-Id"

func (h XAmznTraceIdLoggerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if e.GetCtx() == nil {
		return
	}

	carrier := propagation.MapCarrier{}
	xray.Propagator{}.Inject(e.GetCtx(), &carrier)

	if traceId, ok := carrier[xAmznTraceIdHeader]; ok {
		e.Str(xAmznTraceIdHeader, traceId)
	}
}

const datadogTraceIdKey = "dd.trace_id"
const datadogSpanIdKey = "dd.span_id"

// Extracts trace ID and span ID from the logging context and writes it to the `dd.trace_id` and `dd.span_id` keys in the log event.
//
// Usage:
//
//	log.Logger = log.Logger.Hook(ginruntime.DatadogTraceCorrelationLoggerHook{})
//
// Reference: https://docs.datadoghq.com/logs/log_configuration/pipelines/?tab=traceid#trace-id-attribute
type DatadogTraceCorrelationLoggerHook struct{}

func (h DatadogTraceCorrelationLoggerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if e.GetCtx() == nil {
		return
	}

	spanCtx := trace.SpanContextFromContext(e.GetCtx())

	if spanCtx.IsValid() {
		e.Str(datadogTraceIdKey, spanCtx.TraceID().String())
		e.Str(datadogSpanIdKey, spanCtx.SpanID().String())
	}
}
