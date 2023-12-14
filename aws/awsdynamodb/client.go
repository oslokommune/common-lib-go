package awsdynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/oslokommune/common-lib-go/aws/lambdaruntime"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *dynamodb.Client {
	var cfg aws.Config

	if lambdaruntime.IsRunningAsLambda() {
		cfg, _ = config.LoadDefaultConfig(context.TODO())
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID && region == "eu-north-1" {
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

	// Create an Amazon DynamoDB client.
	dynamodbClient := dynamodb.NewFromConfig(cfg)
	return dynamodbClient
}

type DynamoDBDescribeTableAPI interface {
	DescribeTable(ctx context.Context,
		params *dynamodb.DescribeTableInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

type DynamoDBScanTableAPI interface {
	Scan(ctx context.Context,
		params *dynamodb.ScanInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type DynamoDBUpdateItemApi interface {
	UpdateItem(ctx context.Context,
		params *dynamodb.UpdateItemInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

type DynamoDBQueryTableApi interface {
	Query(ctx context.Context,
		params *dynamodb.QueryInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

func describeTable(c context.Context, api DynamoDBDescribeTableAPI, input *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	return api.DescribeTable(c, input)
}

func scanTable(c context.Context, api DynamoDBScanTableAPI, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return api.Scan(c, input)
}

func updateItem(c context.Context, api DynamoDBUpdateItemApi, input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return api.UpdateItem(c, input)
}

func queryTable(c context.Context, api DynamoDBQueryTableApi, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return api.Query(c, input)
}

func UpdateTableItem(ctx context.Context, tablename string, client DynamoDBUpdateItemApi, key map[string]any, values map[string]any) (*dynamodb.UpdateItemOutput, error) {
	pk, err := attributevalue.MarshalMap(key)
	if err != nil {
		return nil, err
	}

	var upd expression.UpdateBuilder
	for k, v := range values {
		upd = upd.Set(expression.Name(k), expression.Value(v))
	}

	expr, err := expression.NewBuilder().WithUpdate(upd).Build()

	input := dynamodb.UpdateItemInput{
		TableName:                 aws.String(tablename),
		Key:                       pk,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	return updateItem(ctx, client, &input)
}

func QueryTable[T any](ctx context.Context, tablename string, client DynamoDBQueryTableApi, keys map[string]any) []T {
	var keyConditionBuilder expression.KeyConditionBuilder
	for k, v := range keys {
		keyBuilder := expression.Key(k).Equal(expression.Value(v))
		if keyConditionBuilder.IsSet() {
			keyConditionBuilder = keyConditionBuilder.And(keyBuilder)
		} else {
			keyConditionBuilder = keyBuilder
		}
	}

	expr, _ := expression.NewBuilder().WithKeyCondition(keyConditionBuilder).Build()
	input := dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	data, err := queryTable(ctx, client, &input)
	if err != nil {
		log.Error().Err(err).Msg("failed to query dynamodb table")
		return nil
	}

	var recs []T
	err = attributevalue.UnmarshalListOfMaps(data.Items, &recs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal database data.")
		return nil
	}

	log.Printf("Loaded following records: %v", recs)
	return recs
}

func ReadAllTableData[T any](ctx context.Context, tablename string, client DynamoDBScanTableAPI) []T {
	input := dynamodb.ScanInput{
		TableName: aws.String(tablename),
	}

	info, err := scanTable(context.Background(), client, &input)
	if err != nil {
		log.Error().Err(err).Msg("could not read from dynamodb.")
		return nil
	}

	var recs []T
	err = attributevalue.UnmarshalListOfMaps(info.Items, &recs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal database data.")
		return nil
	}

	log.Printf("Loaded following records: %v", recs)

	return recs
}
