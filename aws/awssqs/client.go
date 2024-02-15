package awssqs

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *sqs.Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	if useTracing {
		otelaws.AppendMiddlewares(&cfg.APIOptions)
	}

	// Create an Amazon sqs client.
	client := sqs.NewFromConfig(cfg)
	return client
}

type SqsSendMessageApi interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
}

func publish(ctx context.Context, api SqsSendMessageApi, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return api.SendMessage(ctx, input)
}

func PublishMessage(ctx context.Context, client SqsSendMessageApi, queueName string, message *string) (*sqs.SendMessageOutput, error) {
	// Get URL of queue
	getQueueUrlInput := &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}

	result, err := client.GetQueueUrl(ctx, getQueueUrlInput)
	if err != nil {
		return nil, err
	}

	queueURL := result.QueueUrl
	hash := getMD5Hash(message)

	input := sqs.SendMessageInput{
		MessageBody:            message,
		QueueUrl:               queueURL,
		DelaySeconds:           0,
		MessageDeduplicationId: &hash,
		MessageGroupId:         &hash,
	}

	return publish(ctx, client, &input)
}

func getMD5Hash(text *string) string {
	hasher := md5.New()
	hasher.Write([]byte(*text))
	return hex.EncodeToString(hasher.Sum(nil))
}
