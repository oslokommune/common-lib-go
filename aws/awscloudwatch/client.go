package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	metrictypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/oslokommune/common-lib-go/aws/lambdaruntime"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewLogsClient(useTracing bool) *cloudwatchlogs.Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	if useTracing {
		otelaws.AppendMiddlewares(&cfg.APIOptions)
	}

	// Create an Amazon CloudwatchLogs client
	cloudwatchlogsClient := cloudwatchlogs.NewFromConfig(cfg)
	return cloudwatchlogsClient
}

func NewClient(useTracing bool) *cloudwatch.Client {
	var cfg aws.Config

	if lambdaruntime.IsRunningAsLambda() {
		cfg, _ = config.LoadDefaultConfig(context.TODO())
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == cloudwatch.ServiceID && region == "eu-north-1" {
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

	// Create an Amazon Cloudwatch client.
	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	return cloudwatchClient
}

type DescribeLogStreamsApi interface {
	DescribeLogStreams(ctx context.Context, input *cloudwatchlogs.DescribeLogStreamsInput, optionFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
}

type GetLogEventsApi interface {
	GetLogEvents(ctx context.Context, input *cloudwatchlogs.GetLogEventsInput, optionFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogEventsOutput, error)
}

type GetMetricDataApi interface {
	GetMetricData(ctx context.Context, input *cloudwatch.GetMetricDataInput, optionFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error)
}

func describeLogStreams(ctx context.Context, input *cloudwatchlogs.DescribeLogStreamsInput, client DescribeLogStreamsApi) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	return client.DescribeLogStreams(ctx, input)
}

func getLogEvents(ctx context.Context, input *cloudwatchlogs.GetLogEventsInput, client GetLogEventsApi) (*cloudwatchlogs.GetLogEventsOutput, error) {
	return client.GetLogEvents(ctx, input)
}

func getMetricData(ctx context.Context, input *cloudwatch.GetMetricDataInput, client GetMetricDataApi) (*cloudwatch.GetMetricDataOutput, error) {
	return client.GetMetricData(ctx, input)
}

// FetchLogStreams reads available logstreams from the specified cloudwatch LogGroup and returns them as a slice.
// The list is filtered to only include streams for the specified container and ECS taskArn and only includes
// streams that have been written to during the last 30 minutes
func FetchLogStreams(ctx context.Context, logGroupName string, containerName, taskArn *string, client DescribeLogStreamsApi) ([]types.LogStream, error) {
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
		OrderBy:      types.OrderByLastEventTime,
		Limit:        aws.Int32(10),
		Descending:   aws.Bool(true),
	}

	logStreams, err := describeLogStreams(ctx, input, client)
	if err != nil {
		return nil, err
	}

	thirtyMinutesAgo := time.Now().Add(time.Minute * -30).Unix()

	events := make([]types.LogStream, 0)
	for _, item := range logStreams.LogStreams {
		if *item.LastEventTimestamp > (thirtyMinutesAgo * 1000) {
			if containerName != nil && taskArn != nil {
				if strings.Contains(*item.LogStreamName, *containerName) && strings.Contains(*item.LogStreamName, *taskArn) {
					events = append(events, item)
				}
			} else {
				events = append(events, item)
			}
		}
	}

	return events, nil
}

// FetchCloudwatchLogs reads the content of a log stream and returns a two dimensional slice of log events
func FetchCloudwatchLogs(ctx context.Context, logGroupName, logStreamName string, nextForwardToken *string, interval time.Duration, client GetLogEventsApi) (outputList [][]types.OutputLogEvent, nextToken *string, err error) {
	startTime := time.Now().Add(-interval)

	logEventsInput := cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		StartFromHead: aws.Bool(false),
		StartTime:     aws.Int64(startTime.Unix() * 1000),
	}

	nextToken = nextForwardToken

	var logEventsOutput *cloudwatchlogs.GetLogEventsOutput

	for {
		if nextToken != nil {
			logEventsInput.NextToken = nextToken
			// per GetLogEventsInput documentation, this is required if using a previous aquired nextForwardToken
			logEventsInput.StartFromHead = aws.Bool(true)
		}

		logEventsOutput, err = getLogEvents(ctx, &logEventsInput, client)
		if err != nil {
			return
		}

		if nextToken != nil && *logEventsOutput.NextForwardToken == *nextToken {
			break
		}

		nextToken = logEventsOutput.NextForwardToken
		outputList = append(outputList, logEventsOutput.Events)
	}

	return
}

// FetchCpuAndMemoryUsage queries cloudwatch for ECS service metrics and returns the following metrics:
// - MemoryUtilized
// - MemoryReserved
// - CpuUtilized
// - CpuReserved
func FetchCpuAndMemoryUsage(ctx context.Context, name, clusterName string, client GetMetricDataApi) (memoryUtilized uint32, memoryReserved uint32, cpuUtilized uint32, cpuReserved uint32, err error) {
	startTime := time.Now().Add(-10 * time.Minute)
	endTime := time.Now()
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []metrictypes.MetricDataQuery{
			{
				Id: aws.String("mem_used"),
				MetricStat: &metrictypes.MetricStat{
					Metric: &metrictypes.Metric{
						Namespace:  aws.String("ECS/ContainerInsights"),
						MetricName: aws.String("MemoryUtilized"),
						Dimensions: []metrictypes.Dimension{
							{
								Name:  aws.String("TaskDefinitionFamily"),
								Value: aws.String(name),
							},
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
					},
					Period: aws.Int32(60),
					Stat:   aws.String("Average"),
				},
			},
			{
				Id: aws.String("mem_reserved"),
				MetricStat: &metrictypes.MetricStat{
					Metric: &metrictypes.Metric{
						Namespace:  aws.String("ECS/ContainerInsights"),
						MetricName: aws.String("MemoryReserved"),
						Dimensions: []metrictypes.Dimension{
							{
								Name:  aws.String("TaskDefinitionFamily"),
								Value: aws.String(name),
							},
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
					},
					Period: aws.Int32(60),
					Stat:   aws.String("Average"),
				},
			},
			{
				Id: aws.String("cpu_used"),
				MetricStat: &metrictypes.MetricStat{
					Metric: &metrictypes.Metric{
						Namespace:  aws.String("ECS/ContainerInsights"),
						MetricName: aws.String("CpuUtilized"),
						Dimensions: []metrictypes.Dimension{
							{
								Name:  aws.String("TaskDefinitionFamily"),
								Value: aws.String(name),
							},
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
					},
					Period: aws.Int32(60),
					Stat:   aws.String("Average"),
				},
			},
			{
				Id: aws.String("cpu_reserved"),
				MetricStat: &metrictypes.MetricStat{
					Metric: &metrictypes.Metric{
						Namespace:  aws.String("ECS/ContainerInsights"),
						MetricName: aws.String("CpuReserved"),
						Dimensions: []metrictypes.Dimension{
							{
								Name:  aws.String("TaskDefinitionFamily"),
								Value: aws.String(name),
							},
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
					},
					Period: aws.Int32(60),
					Stat:   aws.String("Average"),
				},
			},
		},
		StartTime: &startTime,
		EndTime:   &endTime,
	}

	// Call CloudWatch's GetMetricData API to retrieve metric data
	result, err := getMetricData(ctx, input, client)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	datapoints := uint32(0)

	// Process the returned data
	for _, metricDataResult := range result.MetricDataResults {
		for _, value := range metricDataResult.Values {
			switch *metricDataResult.Id {
			case "mem_used":
				memoryUtilized += uint32(value)
			case "mem_reserved":
				memoryReserved += uint32(value)
			case "cpu_used":
				cpuUtilized += uint32(value)
			case "cpu_reserved":
				cpuReserved += uint32(value)
				datapoints++
			}
		}
	}

	if datapoints == 0 {
		return 0, 0, 0, 0, nil
	}

	return (cpuUtilized / datapoints), (cpuReserved / datapoints), (memoryUtilized / datapoints), (memoryReserved / datapoints), nil
}
