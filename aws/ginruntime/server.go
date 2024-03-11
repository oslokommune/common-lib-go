package ginruntime

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
)

func (e *GinEngine) lambdaProxy() func(ctx context.Context, req any) (any, error) {
	var proxy any

	proxy = func(ctx context.Context, req any) (any, error) {
		typedReq, ok := req.(events.APIGatewayV2HTTPRequest)
		if !ok {
			return events.APIGatewayV2HTTPResponse{StatusCode: 500}, fmt.Errorf("Unsupported request type: %T", req)
		}
		return ginadapter.NewV2(e.engine).ProxyWithContext(ctx, typedReq)
	}

	if e.TracingEnabled() {
		proxy = otellambda.InstrumentHandler(proxy, e.OTelLambdaOptions()...).(func(context.Context, any) (any, error))
	}
	return proxy.(func(context.Context, any) (any, error))
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
