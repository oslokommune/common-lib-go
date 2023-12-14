package awssecretsmanager

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go"
	"github.com/oslokommune/common-lib-go/aws/lambdaruntime"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *secretsmanager.Client {
	var cfg aws.Config

	if lambdaruntime.IsRunningAsLambda() {
		cfg, _ = config.LoadDefaultConfig(context.TODO())
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			if service == secretsmanager.ServiceID && region == "eu-north-1" {
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

	// Create an Amazon SecretsMananger client.
	client := secretsmanager.NewFromConfig(cfg)
	return client
}

type GetSecretValueApi interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

func getSecretValue(ctx context.Context, api GetSecretValueApi, input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return api.GetSecretValue(ctx, input)
}

// fetches SecretsMananger value with context which enables instrumenting
func GetSecret(ctx context.Context, client GetSecretValueApi, secretName string) (*string, error) {
	input := secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := getSecretValue(ctx, client, &input)

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			log.Printf("call to secretsmanager failed with code: %s, message: %s, fault: %s", ae.ErrorCode(), ae.ErrorMessage(), ae.ErrorFault().String())
		}
		return nil, err
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		return &secretString, nil
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			log.Err(err).Msg("Base64 Decode Error")
			return nil, err
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		return &decodedBinarySecret, nil
	}
}
