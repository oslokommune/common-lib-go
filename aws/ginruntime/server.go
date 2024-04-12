package ginruntime

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
)

func (e *GinEngine) lambdaProxy() any {
	var proxy any

	proxy = func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return ginadapter.NewV2(e.engine).ProxyWithContext(ctx, req)
	}

	if e.TracingEnabled() {
		proxy = otellambda.InstrumentHandler(proxy, e.oTelLambdaOptions()...)
	}

	return proxy
}

func (e *GinEngine) StartServer() {
	defer e.shutdownCallbacks()

	if IsRunningAsLambda() {
		proxy := e.lambdaProxy()
		lambda.StartWithOptions(proxy, lambda.WithContext(e.ctx))
	} else {
		if err := e.engine.Run(); err != nil {
			log.Error().Msgf("Error starting gin %v", err)
		}
	}

	log.Info().Msg("Application exiting.")
}

// Used for unit testing
func (e *GinEngine) ServerHttp(recorder *httptest.ResponseRecorder, request *http.Request) {
	e.engine.ServeHTTP(recorder, request)
}
