package awssns

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/oslokommune/common-lib-go/aws/lambdaruntime"
)

func NewSNSClient() *sns.Client {
	var cfg aws.Config

	if lambdaruntime.IsRunningAsLambda() {
		cfg, _ = config.LoadDefaultConfig(context.TODO())
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			if service == sns.ServiceID && region == "eu-north-1" {
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

	// Create an Amazon SNS client.
	client := sns.NewFromConfig(cfg)
	return client
}

type SNSPublishApi interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

func publish(ctx context.Context, api SNSPublishApi, input *sns.PublishInput) (*sns.PublishOutput, error) {
	return api.Publish(ctx, input)
}

func PublishToTopic(ctx context.Context, client SNSPublishApi, topicArn, message, subject string) (*sns.PublishOutput, error) {
	input := sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicArn),
		Subject:  aws.String(subject),
	}

	return publish(ctx, client, &input)
}
