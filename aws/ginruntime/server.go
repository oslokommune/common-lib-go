package ginruntime

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/oslokommune/common-lib-go/aws/lambdaruntime"
	"github.com/rs/zerolog/log"
)

func (e *GinEngine) StartServer(ctx context.Context) {
	// Check if running as a lambda function
	if lambdaruntime.IsRunningAsLambda() {
		proxy := func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return ginadapter.NewV2(e.engine).ProxyWithContext(ctx, req)
		}
		lambda.StartWithOptions(proxy, lambda.WithContext(ctx))
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
