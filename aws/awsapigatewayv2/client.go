package awsapigatewayv2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	"github.com/oslokommune/common-lib-go/aws/lambdaruntime"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *apigatewayv2.Client {
	var cfg aws.Config

	if lambdaruntime.IsRunningAsLambda() {
		cfg, _ = config.LoadDefaultConfig(context.TODO())
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == apigatewayv2.ServiceID && region == "eu-north-1" {
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
	apigatewayv2Client := apigatewayv2.NewFromConfig(cfg)
	return apigatewayv2Client
}

type ApiGatewayv2GetApiMappingApi interface {
	GetApiMappings(ctx context.Context, params *apigatewayv2.GetApiMappingsInput, optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.GetApiMappingsOutput, error)
}

type ApiGatewayv2GetDomainNames interface {
	GetDomainNames(ctx context.Context, params *apigatewayv2.GetDomainNamesInput, optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.GetDomainNamesOutput, error)
}

type ApiGatewayv2GetApisApi interface {
	GetApis(ctx context.Context,
		params *apigatewayv2.GetApisInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.GetApisOutput, error)
}

func getDomainNames(ctx context.Context, api ApiGatewayv2GetDomainNames, input *apigatewayv2.GetDomainNamesInput) (*apigatewayv2.GetDomainNamesOutput, error) {
	return api.GetDomainNames(ctx, input)
}

func getMappings(ctx context.Context, api ApiGatewayv2GetApiMappingApi, input *apigatewayv2.GetApiMappingsInput) (*apigatewayv2.GetApiMappingsOutput, error) {
	return api.GetApiMappings(ctx, input)
}

func listApis(ctx context.Context, api ApiGatewayv2GetApisApi, input *apigatewayv2.GetApisInput) (*apigatewayv2.GetApisOutput, error) {
	return api.GetApis(ctx, input)
}

func listDomainNames(ctx context.Context, api ApiGatewayv2GetDomainNames, input *apigatewayv2.GetDomainNamesInput) (*apigatewayv2.GetDomainNamesOutput, error) {
	return api.GetDomainNames(ctx, input)
}

func GetApiMappings(ctx context.Context, client ApiGatewayv2GetApiMappingApi, domainName string) ([]types.ApiMapping, error) {
	input := &apigatewayv2.GetApiMappingsInput{
		DomainName: &domainName,
	}

	mappings, err := getMappings(ctx, client, input)
	if err != nil {
		return nil, err
	}

	return mappings.Items, nil
}

func GetDomainNames(ctx context.Context, client ApiGatewayv2GetDomainNames) ([]types.DomainName, error) {
	input := &apigatewayv2.GetDomainNamesInput{}

	list, err := listDomainNames(ctx, client, input)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetApiGatewayv2Apis(ctx context.Context, client ApiGatewayv2GetApisApi) ([]types.Api, error) {
	input := &apigatewayv2.GetApisInput{
		MaxResults: aws.String("10"),
	}

	list, err := listApis(ctx, client, input)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
