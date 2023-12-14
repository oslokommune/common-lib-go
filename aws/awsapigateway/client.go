package awsapigateway

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/oslokommune/common-lib-go/aws/lambdaruntime"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *apigateway.Client {
	var cfg aws.Config

	if lambdaruntime.IsRunningAsLambda() {
		cfg, _ = config.LoadDefaultConfig(context.TODO())
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == apigateway.ServiceID && region == "eu-north-1" {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://localhost:4566",
					SigningRegion: "eu-north-1",
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})

		// Use the SDK's default configuration with region and custome endpoint resolver
		cfg, _ = config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"), config.WithEndpointResolverWithOptions(customResolver))
	}

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
