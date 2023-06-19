package awss3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/oslokommune/common-lib-go/lambdaruntime"

	"github.com/rs/zerolog/log"
)

type S3File struct {
	Name      string
	CreatedAt time.Time
}

func NewS3Client() *s3.Client {
	var cfg aws.Config

	if lambdaruntime.IsRunningAsLambda() {
		cfg, _ = config.LoadDefaultConfig(context.TODO())
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			if service == s3.ServiceID && region == "eu-north-1" {
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

	s3client := s3.NewFromConfig(cfg)
	return s3client
}

type ListObjectsV2API interface {
	ListObjectsV2(ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func listObjects(ctx context.Context, api ListObjectsV2API, params *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return api.ListObjectsV2(ctx, params)
}

func ListBucketObjects(ctx context.Context, client ListObjectsV2API, bucketName string, prefix string) ([]S3File, error) {
	input := s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	output, err := listObjects(ctx, client, &input)
	if err != nil {
		log.Error().Err(err).Msg("failed to read bucket content")
		return nil, err
	}

	list := make([]S3File, 0, output.KeyCount)
	for _, v := range output.Contents {
		list = append(list, S3File{
			Name:      *v.Key,
			CreatedAt: *v.LastModified,
		})
	}

	return list, nil
}
