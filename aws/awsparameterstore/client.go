package awsssm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/oslokommune/common-lib-go/lambdaruntime"
)

func NewSSMClient(ctx context.Context) *ssm.Client {
	var cfg aws.Config

	if lambdaruntime.IsRunningAsLambda() {
		cfg, _ = config.LoadDefaultConfig(ctx)
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			if service == ssm.ServiceID && region == "eu-north-1" {
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

	// Create an Amazon SecretsMananger client.
	client := ssm.NewFromConfig(cfg)
	return client
}

type GetParameterAPI interface {
	GetParameter(ctx context.Context,
		params *ssm.GetParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

func getParameter(ctx context.Context, api GetParameterAPI, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	return api.GetParameter(ctx, input)
}

func GetParameterStoreParameter(ctx context.Context, client GetParameterAPI, name string, container any) error {
	bool := aws.Bool(true)
	input := ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: bool,
	}

	output, err := getParameter(ctx, client, &input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(*output.Parameter.Value), container); err != nil {
		return err
	}

	return nil
}
