package ginruntime

import (
	"github.com/oslokommune/common-lib-go/aws/ginruntime/openapi"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type OpenAPIOptions struct {
	service          string
	version          string
	description      string
	swaggerUiDistUrl string
}

type TracingOptions struct {
	service    string
	tp         *trace.TracerProvider
	propagator propagation.TextMapPropagator
}

type Option struct {
	openapi *OpenAPIOptions
	tracing *TracingOptions
}

// Enables OpenAPI endpoint `/openapi.json` and Swagger UI endpoint `/docs`.
//
// Routes must be annotated with `openapi.Annotate` to be included in the OpenAPI spec.
//
// The `/docs` endpoint uses `swaggerUiDistUrl` to load JavaScript and CSS for Swagger UI.
// See here for more information: https://github.com/swagger-api/swagger-ui/blob/master/docs/usage/installation.md
func WithOpenAPI(
	service string,
	version string,
	description string,
	swaggerUiDistUrl string,
) Option {
	return Option{
		openapi: &OpenAPIOptions{
			service:          service,
			version:          version,
			description:      description,
			swaggerUiDistUrl: swaggerUiDistUrl,
		},
	}
}

func (e *GinEngine) OpenAPIEnabled() bool {
	return e.openapi != nil
}

// Configures OpenTelemetry for tracing that integrates with AWS X-Ray and adds a middleware that traces incoming requests.
func WithXRayTracing(
	service string,
) Option {
	return Option{
		tracing: &TracingOptions{service, nil, nil},
	}
}

// Configures OpenTelemetry for tracing and adds a middleware that traces incoming requests.
func WithTracing(
	service string,
	tp *trace.TracerProvider,
	propagator propagation.TextMapPropagator,
) Option {
	return Option{
		tracing: &TracingOptions{service, tp, propagator},
	}
}

func (e *GinEngine) apply(options ...Option) {
	for _, option := range options {
		if option.tracing != nil {
			log.Info().Msg("Enabling tracing")
			e.enableTracing(option.tracing)
		}
	}

	for _, option := range options {
		if option.openapi != nil {
			log.Info().Msgf("Enabling OpenAPI serving static files from %s", option.openapi.swaggerUiDistUrl)
			e.enableOpenAPI(option.openapi)
		}
	}
}

func (e *GinEngine) TracingEnabled() bool {
	return e.tp != nil
}

// Configures OpenTelemetry for tracing.
// Must be called before adding routes to the engine.
func (e *GinEngine) enableTracing(options *TracingOptions) {
	var err error

	if 0 < len(e.engine.Routes()) {
		log.Fatal().Msg("tracing must be enabled before adding routes to the engine, otherwise the middleware isn't applied to the routes")
	}

	if options.service == "" {
		log.Warn().Msg("missing service name to use for tracing")
	}

	if options.tp == nil {
		options.tp, err = xrayconfig.NewTracerProvider(e.ctx)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to initialize X-Ray TracerProvider")
		}
	}

	if options.propagator == nil {
		options.propagator = &xray.Propagator{}
	}

	e.tp = options.tp
	e.propagator = options.propagator

	e.configureOpenTelemetry()
	e.useHttpServerSpanMiddleware(options.service)

	e.OnShutdown(func() {
		err := e.tp.Shutdown(e.ctx)
		if err != nil {
			log.Error().Err(err).Msg("Error shutting down tracer provider")
		}
	})
}

func (e *GinEngine) enableOpenAPI(options *OpenAPIOptions) {
	e.openapi = openapi.New(options.service, options.version, options.description, options.swaggerUiDistUrl)

	e.AddRoute(nil, "/openapi.json", GET, nil, e.openapi.JsonSpecRoute)
	e.AddRoute(nil, "/docs", GET, nil, e.openapi.UiRoute)
}
