package awslambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/oslokommune/common-lib-go/aws/lambdaruntime"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *lambda.Client {
	var cfg aws.Config

	if lambdaruntime.IsRunningAsLambda() {
		cfg, _ = config.LoadDefaultConfig(context.TODO())
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			if service == lambda.ServiceID && region == "eu-north-1" {
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

	lambdaClient := lambda.NewFromConfig(cfg)
	return lambdaClient
}

type UpdateFunctionCodeApi interface {
	UpdateFunctionCode(ctx context.Context,
		params *lambda.UpdateFunctionCodeInput,
		optFns ...func(*lambda.Options)) (*lambda.UpdateFunctionCodeOutput, error)
}

type GetFunctionApi interface {
	GetFunction(ctx context.Context,
		params *lambda.GetFunctionInput,
		optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error)
}

type TagResourceApi interface {
	TagResource(ctx context.Context,
		params *lambda.TagResourceInput,
		optFns ...func(*lambda.Options)) (*lambda.TagResourceOutput, error)
}

type UpdateFunctionConfigurationApi interface {
	UpdateFunctionConfiguration(ctx context.Context,
		params *lambda.UpdateFunctionConfigurationInput,
		optFns ...func(*lambda.Options)) (*lambda.UpdateFunctionConfigurationOutput, error)
}

type GetFunctionConfigurationApi interface {
	GetFunctionConfiguration(ctx context.Context, params *lambda.GetFunctionConfigurationInput,
		optFns ...func(*lambda.Options)) (*lambda.GetFunctionConfigurationOutput, error)
}

type ListFunctionsApi interface {
	ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput,
		optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}

func listFunctions(ctx context.Context, api ListFunctionsApi, input *lambda.ListFunctionsInput) (*lambda.ListFunctionsOutput, error) {
	return api.ListFunctions(ctx, input)
}

func getFunction(ctx context.Context, api GetFunctionApi, input *lambda.GetFunctionInput) (*lambda.GetFunctionOutput, error) {
	return api.GetFunction(ctx, input)
}

func updateFunctionCode(ctx context.Context, api UpdateFunctionCodeApi, input *lambda.UpdateFunctionCodeInput) (*lambda.UpdateFunctionCodeOutput, error) {
	return api.UpdateFunctionCode(ctx, input)
}

func tagResource(ctx context.Context, api TagResourceApi, input *lambda.TagResourceInput) (*lambda.TagResourceOutput, error) {
	return api.TagResource(ctx, input)
}

func updateFunctionConfiguration(ctx context.Context, api UpdateFunctionConfigurationApi, input *lambda.UpdateFunctionConfigurationInput) (*lambda.UpdateFunctionConfigurationOutput, error) {
	return api.UpdateFunctionConfiguration(ctx, input)
}

func getFunctionConfiguration(ctx context.Context, api GetFunctionConfigurationApi, input *lambda.GetFunctionConfigurationInput) (*lambda.GetFunctionConfigurationOutput, error) {
	return api.GetFunctionConfiguration(ctx, input)
}

// List all lambda ListFunctions
func ListFunctions(ctx context.Context, client ListFunctionsApi) (*lambda.ListFunctionsOutput, error) {
	input := &lambda.ListFunctionsInput{
		MaxItems: aws.Int32(100),
	}

	return listFunctions(ctx, client, input)
}

// Gets lambda function details by function name
func GetFunction(ctx context.Context, client GetFunctionApi, functionName string) (*lambda.GetFunctionOutput, error) {
	getFunctionInput := &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	}

	output, err := getFunction(ctx, client, getFunctionInput)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// Updates lambda function code referenced by function name to code uploaded to specified s3 bucket
func UpdateFunctionCode(ctx context.Context, client UpdateFunctionCodeApi, functionName string, s3Bucket string, s3Key string, architecture types.Architecture) (*lambda.UpdateFunctionCodeOutput, error) {
	log.Info().Msgf("updating lambda function with name %s from bucket %s and file %s", functionName, s3Bucket, s3Key)

	updateFunctionInput := &lambda.UpdateFunctionCodeInput{
		FunctionName:  aws.String(functionName),
		S3Bucket:      aws.String(s3Bucket),
		S3Key:         aws.String(fmt.Sprintf("%s/%s", functionName, s3Key)),
		Architectures: []types.Architecture{architecture},
	}

	output, err := updateFunctionCode(ctx, client, updateFunctionInput)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func TagLambdaFunction(ctx context.Context, client TagResourceApi, functionName string, tags map[string]string) error {
	tagResourceInput := &lambda.TagResourceInput{
		Resource: aws.String(functionName),
		Tags:     tags,
	}

	_, err := tagResource(ctx, client, tagResourceInput)
	return err
}

func GetFunctionDescription(ctx context.Context, client GetFunctionConfigurationApi, functionName string) (*string, error) {
	input := lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
	}

	config, err := getFunctionConfiguration(ctx, client, &input)
	if err != nil {
		return nil, err
	}

	return config.Description, nil
}

func UpdateFunctionDescription(ctx context.Context, client UpdateFunctionConfigurationApi, functionName string, description string) error {
	input := lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
		Description:  aws.String(description),
	}

	_, err := updateFunctionConfiguration(ctx, client, &input)
	return err
}
