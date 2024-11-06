package awsssm

import (
	"context"
	"encoding/json"

	//"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *ssm.Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	if useTracing {
		otelaws.AppendMiddlewares(&cfg.APIOptions)
	}

	// Create an Amazon SystemsManager client.
	client := ssm.NewFromConfig(cfg)
	return client
}

type GetParameterAPI interface {
	GetParameter(ctx context.Context,
		params *ssm.GetParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

type PutParameterApi interface {
	PutParameter(ctx context.Context,
		params *ssm.PutParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
}

type DescribeParametersApi interface {
	DescribeParameters(ctx context.Context,
		params *ssm.DescribeParametersInput,
		optFns ...func(*ssm.Options)) (*ssm.DescribeParametersOutput, error)
}

func getParameter(ctx context.Context, api GetParameterAPI, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	return api.GetParameter(ctx, input)
}

func putParameter(ctx context.Context, api PutParameterApi, input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	return api.PutParameter(ctx, input)
}

func describeParameters(ctx context.Context, api DescribeParametersApi, input *ssm.DescribeParametersInput) (*ssm.DescribeParametersOutput, error) {
	return api.DescribeParameters(ctx, input)
}

func GetParameterStoreParameterString(ctx context.Context, client GetParameterAPI, name string) (*string, error) {
	bool := aws.Bool(true)
	input := ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: bool,
	}

	output, err := getParameter(ctx, client, &input)
	if err != nil {
		return nil, err
	}

	return output.Parameter.Value, nil
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

func UpdateParameterStoreParameter(ctx context.Context, client PutParameterApi, name, value string) error {
	overwrite := aws.Bool(true)
	input := ssm.PutParameterInput{
		Name:      &name,
		Value:     &value,
		Overwrite: overwrite,
	}

	_, err := putParameter(ctx, client, &input)

	return err
}

func DescribeParameterStoreParameter(ctx context.Context, client DescribeParametersApi, name string) (*ssm.DescribeParametersOutput, error) {
	filter := types.ParameterStringFilter{
		Key:    aws.String("Name"),
		Values: []string{name},
	}

	input := ssm.DescribeParametersInput{
		ParameterFilters: []types.ParameterStringFilter{filter},
	}

	return describeParameters(ctx, client, &input)
}
