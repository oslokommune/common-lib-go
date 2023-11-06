package lambdaruntime

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitializeTracing(ctx context.Context) *trace.TracerProvider {
	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		log.Panic().Err(err).Msg("Error creating trace provider")
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return tp
}

func ShutdownTracing(ctx context.Context, tp *trace.TracerProvider) {
	err := tp.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error shutting down tracer provider")
	}
}
