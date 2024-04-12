package ginruntime

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestServerLambdaProxyWithoutTracing(t *testing.T) {
	engine := New(context.Background())
	proxy := engine.lambdaProxy().(func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error))
	event := events.APIGatewayV2HTTPRequest{}
	proxy(context.Background(), event)
}

func TestServerLambdaProxyWithTracing(t *testing.T) {
	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test")
	lc := lambdacontext.LambdaContext{AwsRequestID: "test", InvokedFunctionArn: "test", Identity: lambdacontext.CognitoIdentity{}, ClientContext: lambdacontext.ClientContext{}}
	ctx := lambdacontext.NewContext(context.Background(), &lc)
	engine := New(ctx, WithTracing("test", NewInterceptingTracerProvider(func(span []trace.ReadOnlySpan) {}), &xray.Propagator{}))
	proxy := engine.lambdaProxy().(func(context.Context, any) (any, error))
	event := events.APIGatewayV2HTTPRequest{}
	proxy(ctx, event)
}
