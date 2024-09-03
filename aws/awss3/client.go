package awss3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

type S3File struct {
	CreatedAt time.Time
	Name      string
}

func NewClient(useTracing bool) *s3.Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	if useTracing {
		otelaws.AppendMiddlewares(&cfg.APIOptions)
	}

	s3client := s3.NewFromConfig(cfg)
	return s3client
}

type ListObjectsV2API interface {
	ListObjectsV2(ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type GetObjectAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func listObjects(ctx context.Context, api ListObjectsV2API, params *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return api.ListObjectsV2(ctx, params)
}

func getObject(ctx context.Context, api GetObjectAPI, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return api.GetObject(ctx, params)
}

func getObjectWithManager(ctx context.Context, api GetObjectAPI, params *s3.GetObjectInput, buffer *manager.WriteAtBuffer) error {
	var partMiBs int64 = 10
	downloader := manager.NewDownloader(api, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024
	})

	_, err := downloader.Download(ctx, buffer, params)
	if err != nil {
		return err
	}

	return err
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

	list := make([]S3File, 0, *output.KeyCount)
	for _, v := range output.Contents {
		list = append(list, S3File{
			Name:      *v.Key,
			CreatedAt: *v.LastModified,
		})
	}

	return list, nil
}

// DownloadFile downloads file from S3 and returns it as a byte slice
func DownloadFile(ctx context.Context, api GetObjectAPI, bucketName string, objectKey string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	output, err := getObject(ctx, api, input)
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	bytes, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func GeneratePresignedURL(ctx context.Context, api *s3.Client, bucketName string, objectKey string, mimeType string, fileName string, expiry time.Duration) (string, error) {
	presigner := s3.NewPresignClient(api)

	input := &s3.GetObjectInput{
		Bucket:                     aws.String(bucketName),
		Key:                        aws.String(objectKey),
		ResponseContentType:        aws.String(mimeType),
		ResponseContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"%s\"", fileName)),
	}

	// Generer en presigned URL
	presignedURL, err := presigner.PresignGetObject(ctx, input, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.URL, nil
}

// DownloadFileStream downloads file from S3 and returns the io.ReadCloser. This must be closed by the callee function!
func DownloadFileStream(ctx context.Context, api GetObjectAPI, bucketName string, objectKey string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	output, err := getObject(ctx, api, input)
	if err != nil {
		return nil, err
	}

	return output.Body, nil
}

// DownloadFileLarge downloads file from S3 using download manager and returns it as a byte slice
func DownloadLargeFile(ctx context.Context, api GetObjectAPI, bucketName string, objectKey string) ([]byte, error) {
	buffer := manager.NewWriteAtBuffer([]byte{})

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	err := getObjectWithManager(ctx, api, input, buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
