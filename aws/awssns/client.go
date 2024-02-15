package awssns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *sns.Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	if useTracing {
		otelaws.AppendMiddlewares(&cfg.APIOptions)
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
