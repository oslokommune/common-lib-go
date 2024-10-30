package awsecr

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *ecr.Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	if useTracing {
		otelaws.AppendMiddlewares(&cfg.APIOptions)
	}

	// Create an Amazon DynamoDB client.
	ecrClient := ecr.NewFromConfig(cfg)
	return ecrClient
}

type DescribeImagesAPI interface {
	DescribeImages(ctx context.Context,
		params *ecr.DescribeImagesInput,
		optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error)
}

func describeImages(ctx context.Context, api DescribeImagesAPI, input *ecr.DescribeImagesInput) (*ecr.DescribeImagesOutput, error) {
	return api.DescribeImages(ctx, input)
}

func DescribeImages(ctx context.Context, repositoryName string, client DescribeImagesAPI) (*ecr.DescribeImagesOutput, error) {
	maxResults := int32(10)

	input := ecr.DescribeImagesInput{
		RepositoryName: &repositoryName,
		MaxResults:     &maxResults,
	}

	return describeImages(ctx, client, &input)
}
