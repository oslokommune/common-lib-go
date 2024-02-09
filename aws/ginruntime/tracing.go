package ginruntime

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func (e *GinEngine) TracingEnabled() bool {
	return e.tp != nil
}

// Configures OpenTelemetry for tracing that integrates with AWS X-Ray.
// See `EnableCustomTracing` for more info.
func (e *GinEngine) EnableTracing(service string) {
	if service == "" {
		log.Fatal().Msg("missing service name to use for tracing")
	}

	tp, err := xrayconfig.NewTracerProvider(e.ctx)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to initialize X-Ray TracerProvider")
	}
	propagator := &xray.Propagator{}
	e.EnableCustomTracing(service, tp, propagator)
}

// Configures OpenTelemetry for tracing.
// Must be called before adding routes to the engine.
func (e *GinEngine) EnableCustomTracing(service string, tp *trace.TracerProvider, propagator propagation.TextMapPropagator) {
	if 0 < len(e.engine.Routes()) {
		log.Fatal().Msg("tracing must be enabled before adding routes to the engine")
	}

	if tp == nil || propagator == nil {
		log.Fatal().Msg("missing tracer provider or propagator")
	}

	e.tp = tp
	e.propagator = propagator

	e.configureOpenTelemetry()
	e.useHttpServerSpanMiddleware(service)

	e.OnShutdown(func() {
		err := tp.Shutdown(e.ctx)
		if err != nil {
			log.Error().Err(err).Msg("Error shutting down tracer provider")
		}
	})
}

func (e *GinEngine) OTelLambdaOptions() []otellambda.Option {
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
