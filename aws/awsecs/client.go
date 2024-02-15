package awsecs

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(useTracing bool) *ecs.Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	if useTracing {
		otelaws.AppendMiddlewares(&cfg.APIOptions)
	}

	// Create an AWS ECS client.
	ecsClient := ecs.NewFromConfig(cfg)
	return ecsClient
}

type DescribeServicesApi interface {
	DescribeServices(ctx context.Context, params *ecs.DescribeServicesInput, optFns ...func(*ecs.Options)) (*ecs.DescribeServicesOutput, error)
}

type UpdateServiceApi interface {
	UpdateService(ctx context.Context, params *ecs.UpdateServiceInput, optFns ...func(*ecs.Options)) (*ecs.UpdateServiceOutput, error)
}

type TagResourceApi interface {
	TagResource(ctx context.Context, params *ecs.TagResourceInput, optFns ...func(*ecs.Options)) (*ecs.TagResourceOutput, error)
}

type DescribeTaskDefinitionApi interface {
	DescribeTaskDefinition(ctx context.Context, params *ecs.DescribeTaskDefinitionInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTaskDefinitionOutput, error)
}

type RegisterTaskDefinitionApi interface {
	RegisterTaskDefinition(ctx context.Context, params *ecs.RegisterTaskDefinitionInput, optFns ...func(*ecs.Options)) (*ecs.RegisterTaskDefinitionOutput, error)
}

type StopTaskApi interface {
	StopTask(ctx context.Context, params *ecs.StopTaskInput, optFns ...func(*ecs.Options)) (*ecs.StopTaskOutput, error)
}

type ListTasksApi interface {
	ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error)
}

type DescribeClustersApi interface {
	DescribeClusters(ctx context.Context, params *ecs.DescribeClustersInput, optFns ...func(*ecs.Options)) (*ecs.DescribeClustersOutput, error)
}

type ListServicesApi interface {
	ListServices(ctx context.Context, params *ecs.ListServicesInput, optFns ...func(*ecs.Options)) (*ecs.ListServicesOutput, error)
}

type DescribeTasksApi interface {
	DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)
}

type ListContainerInstancesApi interface {
	ListContainerInstances(ctx context.Context, params *ecs.ListContainerInstancesInput, optFns ...func(*ecs.Options)) (*ecs.ListContainerInstancesOutput, error)
}

type DescribeContainerInstancesApi interface {
	DescribeContainerInstances(ctx context.Context, params *ecs.DescribeContainerInstancesInput, optFns ...func(*ecs.Options)) (*ecs.DescribeContainerInstancesOutput, error)
}

type ListClustersApi interface {
	ListClusters(ctx context.Context, params *ecs.ListClustersInput, optFns ...func(*ecs.Options)) (*ecs.ListClustersOutput, error)
}

type ECSServiceApi interface {
	DescribeServicesApi
	UpdateServiceApi
	TagResourceApi
	DescribeTaskDefinitionApi
	RegisterTaskDefinitionApi
	StopTaskApi
	ListTasksApi
	DescribeClustersApi
	ListServicesApi
	DescribeTasksApi
	ListContainerInstancesApi
	DescribeContainerInstancesApi
	ListClustersApi
}

func listClusters(ctx context.Context, api ListClustersApi, input *ecs.ListClustersInput) (*ecs.ListClustersOutput, error) {
	return api.ListClusters(ctx, input)
}

func listServices(ctx context.Context, api ListServicesApi, input *ecs.ListServicesInput) (*ecs.ListServicesOutput, error) {
	return api.ListServices(ctx, input)
}

func describeServices(ctx context.Context, api DescribeServicesApi, input *ecs.DescribeServicesInput) (*ecs.DescribeServicesOutput, error) {
	return api.DescribeServices(ctx, input)
}

func updateService(ctx context.Context, api UpdateServiceApi, input *ecs.UpdateServiceInput) (*ecs.UpdateServiceOutput, error) {
	return api.UpdateService(ctx, input)
}

func tagResource(ctx context.Context, api TagResourceApi, input *ecs.TagResourceInput) (*ecs.TagResourceOutput, error) {
	return api.TagResource(ctx, input)
}

func describeTaskDefinition(ctx context.Context, api DescribeTaskDefinitionApi, input *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	return api.DescribeTaskDefinition(ctx, input)
}

func registerTaskDefinition(ctx context.Context, api RegisterTaskDefinitionApi, input *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	return api.RegisterTaskDefinition(ctx, input)
}

func listTasks(ctx context.Context, api ListTasksApi, input *ecs.ListTasksInput) (*ecs.ListTasksOutput, error) {
	return api.ListTasks(ctx, input)
}

func describeClusters(ctx context.Context, api DescribeClustersApi, input *ecs.DescribeClustersInput) (*ecs.DescribeClustersOutput, error) {
	return api.DescribeClusters(ctx, input)
}

func describeTasks(ctx context.Context, api DescribeTasksApi, input *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error) {
	return api.DescribeTasks(ctx, input)
}

func listContainerInstances(ctx context.Context, api ListContainerInstancesApi, input *ecs.ListContainerInstancesInput) (*ecs.ListContainerInstancesOutput, error) {
	return api.ListContainerInstances(ctx, input)
}

func describeContainerInstances(ctx context.Context, api DescribeContainerInstancesApi, input *ecs.DescribeContainerInstancesInput) (*ecs.DescribeContainerInstancesOutput, error) {
	return api.DescribeContainerInstances(ctx, input)
}

func ListClusters(ctx context.Context, client ListClustersApi) (*ecs.ListClustersOutput, error) {
	input := &ecs.ListClustersInput{}

	return listClusters(ctx, client, input)
}

func ListContainerInstances(ctx context.Context, client ListContainerInstancesApi, clusterName string) (*ecs.ListContainerInstancesOutput, error) {
	input := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(clusterName),
	}

	return listContainerInstances(ctx, client, input)
}

func DescribeContainerInstances(ctx context.Context, client DescribeContainerInstancesApi, clusterName string, containerInstanceArns []string) (*ecs.DescribeContainerInstancesOutput, error) {
	input := &ecs.DescribeContainerInstancesInput{
		ContainerInstances: containerInstanceArns,
		Cluster:            aws.String(clusterName),
	}

	return describeContainerInstances(ctx, client, input)
}

func DescribeTasks(ctx context.Context, client DescribeTasksApi, taskArns []string, clusterName string) (*ecs.DescribeTasksOutput, error) {
	input := &ecs.DescribeTasksInput{
		Tasks:   taskArns,
		Cluster: aws.String(clusterName),
	}
	return describeTasks(ctx, client, input)
}

func ListTasks(ctx context.Context, client ListTasksApi, clusterName, service string) (*ecs.ListTasksOutput, error) {
	input := &ecs.ListTasksInput{
		Cluster:     &clusterName,
		ServiceName: &service,
	}

	return listTasks(ctx, client, input)
}

func ListServices(ctx context.Context, api ListServicesApi, clusterName string) (*ecs.ListServicesOutput, error) {
	input := &ecs.ListServicesInput{
		Cluster:    aws.String(clusterName),
		MaxResults: aws.Int32(100),
	}

	return api.ListServices(ctx, input)
}

func DescribeClusters(ctx context.Context, client DescribeClustersApi, arns []string, includes []types.ClusterField) (*ecs.DescribeClustersOutput, error) {
	input := &ecs.DescribeClustersInput{
		Clusters: arns,
		Include:  includes,
	}

	return describeClusters(ctx, client, input)
}

func TagResource(ctx context.Context, client TagResourceApi, resource string, tagmap map[string]string) error {
	var tags []types.Tag
	for k, v := range tagmap {
		tags = append(tags, types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	tagResourceInput := &ecs.TagResourceInput{
		ResourceArn: aws.String(resource),
		Tags:        tags,
	}

	_, err := tagResource(ctx, client, tagResourceInput)
	return err
}

func DescribeEcsService(ctx context.Context, client DescribeServicesApi, serviceName string, cluster string) (*ecs.DescribeServicesOutput, error) {
	describeServiceInput := ecs.DescribeServicesInput{
		Services: []string{serviceName},
		Include:  []types.ServiceField{types.ServiceFieldTags},
		Cluster:  aws.String(cluster),
	}

	return describeServices(ctx, client, &describeServiceInput)
}

func DescribeServices(ctx context.Context, client DescribeServicesApi, serviceNames []string, cluster string) (*ecs.DescribeServicesOutput, error) {
	describeServiceInput := ecs.DescribeServicesInput{
		Services: serviceNames,
		Include:  []types.ServiceField{types.ServiceFieldTags},
		Cluster:  aws.String(cluster),
	}

	return describeServices(ctx, client, &describeServiceInput)
}

func DescribeTaskDefinition(ctx context.Context, client DescribeTaskDefinitionApi, taskDefinitionArn string) (*ecs.DescribeTaskDefinitionOutput, error) {
	describeTaskDefinitionInput := ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionArn),
	}

	return describeTaskDefinition(ctx, client, &describeTaskDefinitionInput)
}

func UpdateEcsService(ctx context.Context, client ECSServiceApi, image string, serviceName string, cluster string) (*ecs.UpdateServiceOutput, error) {
	// Find task definition used
	serviceDescription, err := DescribeEcsService(ctx, client, serviceName, cluster)
	if err != nil {
		log.Error().Err(err).Msg("failed to describe ecs service")
		return nil, err
	}

	taskDefinitionArn := *serviceDescription.Services[0].TaskDefinition

	// Read task definition json
	taskDefinition, err := DescribeTaskDefinition(ctx, client, taskDefinitionArn)
	if err != nil {
		log.Error().Err(err).Msgf("failed to read task definition with arn %s", taskDefinitionArn)
		return nil, err
	}

	index := -1
	for i, v := range taskDefinition.TaskDefinition.ContainerDefinitions {
		if strings.Contains(*v.Image, serviceName) {
			index = i
			break
		}
	}

	if index < 0 {
		log.Error().Msgf("failed to find container definition with image that contains string %s", image)
		return nil, fmt.Errorf("failed to find container definition with image that contains string %s", image)
	}

	// updates image in task definition
	taskDefinition.TaskDefinition.ContainerDefinitions[index].Image = aws.String(image)

	// Register new task definition
	registerTaskDefinitionInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    taskDefinition.TaskDefinition.ContainerDefinitions,
		Family:                  taskDefinition.TaskDefinition.Family,
		Cpu:                     taskDefinition.TaskDefinition.Cpu,
		EphemeralStorage:        taskDefinition.TaskDefinition.EphemeralStorage,
		ExecutionRoleArn:        taskDefinition.TaskDefinition.ExecutionRoleArn,
		InferenceAccelerators:   taskDefinition.TaskDefinition.InferenceAccelerators,
		IpcMode:                 taskDefinition.TaskDefinition.IpcMode,
		Memory:                  taskDefinition.TaskDefinition.Memory,
		NetworkMode:             taskDefinition.TaskDefinition.NetworkMode,
		PidMode:                 taskDefinition.TaskDefinition.PidMode,
		PlacementConstraints:    taskDefinition.TaskDefinition.PlacementConstraints,
		ProxyConfiguration:      taskDefinition.TaskDefinition.ProxyConfiguration,
		RequiresCompatibilities: taskDefinition.TaskDefinition.RequiresCompatibilities,
		RuntimePlatform:         taskDefinition.TaskDefinition.RuntimePlatform,
		//		Tags:                    taskDefinition.Tags,
		TaskRoleArn: taskDefinition.TaskDefinition.TaskRoleArn,
		Volumes:     taskDefinition.TaskDefinition.Volumes,
	}

	taskDefinitionOutput, err := registerTaskDefinition(ctx, client, &registerTaskDefinitionInput)
	if err != nil {
		log.Error().Err(err).Msg("failed to register new task defintion")
		return nil, err
	}

	// Updates service to use the newly registered task definition
	updateServiceInput := ecs.UpdateServiceInput{
		TaskDefinition: aws.String(*taskDefinitionOutput.TaskDefinition.TaskDefinitionArn),
		Service:        aws.String(serviceName),
		Cluster:        aws.String(cluster),
	}

	updateServiceOutput, err := updateService(ctx, client, &updateServiceInput)
	if err != nil {
		log.Error().Err(err).Msg("failed to update service")
		return nil, err
	}

	return updateServiceOutput, nil
}

func StopEcsService(ctx context.Context, client ECSServiceApi, name string, clusterName string) error {
	log.Debug().Msgf("trying to restart service %s", name)

	updateServiceInput := ecs.UpdateServiceInput{
		Service:            aws.String(name),
		Cluster:            aws.String(clusterName),
		ForceNewDeployment: true,
	}

	_, err := updateService(ctx, client, &updateServiceInput)
	if err != nil {
		log.Error().Err(err).Msg("failed to update service with ForceNewDeployment turned on")
		return err
	}

	return nil
}
