package awsapigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *apigateway.Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	if useTracing {
		otelaws.AppendMiddlewares(&cfg.APIOptions)
	}

	// Create an Amazon ApiGatewayv2 client.
	apigatewayv2Client := apigateway.NewFromConfig(cfg)
	return apigatewayv2Client
}

type ApiGatewayGetRestApisApi interface {
	GetRestApis(ctx context.Context,
		params *apigateway.GetRestApisInput,
		optFns ...func(*apigateway.Options)) (*apigateway.GetRestApisOutput, error)
}

func listApis(ctx context.Context, api ApiGatewayGetRestApisApi, input *apigateway.GetRestApisInput) (*apigateway.GetRestApisOutput, error) {
	return api.GetRestApis(ctx, input)
}

func GetApiGatewayApis(ctx context.Context, client ApiGatewayGetRestApisApi) ([]types.RestApi, error) {
	input := &apigateway.GetRestApisInput{}

	list, err := listApis(ctx, client, input)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
