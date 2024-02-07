package ginruntime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/oslokommune/common-lib-go/aws/lambdaruntime"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
)

func (e *GinEngine) StartServer(ctx context.Context, tracing bool) {
	if lambdaruntime.IsRunningAsLambda() {
		proxy := func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return ginadapter.NewV2(e.engine).ProxyWithContext(ctx, req)
		}

		if tracing {
			tp, err := xrayconfig.NewTracerProvider(ctx)
			if err != nil {
				log.Panic().Err(err).Msg("Error creating trace provider")
			}

			otel.SetTracerProvider(tp)
			otel.SetTextMapPropagator(xray.Propagator{})
			e.engine.Use(otelgin.Middleware(os.Getenv("APP_LABEL")))

			defer func(ctx context.Context) {
				err := tp.Shutdown(ctx)
				if err != nil {
					log.Error().Err(err).Msg("Error shutting down tracer provider")
				}
			}(ctx)

			lambda.StartWithOptions(otellambda.InstrumentHandler(proxy, xrayconfig.WithRecommendedOptions(tp)...), lambda.WithContext(ctx))
		} else {
			lambda.StartWithOptions(proxy, lambda.WithContext(ctx))
		}
	} else {
		if err := e.engine.Run(); err != nil {
			log.Info().Msgf("Error starting gin %v", err)
		}
		log.Info().Msg("Application exiting.")
	}
}

// Used for unit testing
func (e *GinEngine) ServerHttp(recorder *httptest.ResponseRecorder, request *http.Request) {
	e.engine.ServeHTTP(recorder, request)
}
