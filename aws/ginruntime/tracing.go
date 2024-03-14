package ginruntime

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
