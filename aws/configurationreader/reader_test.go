package configurationreader

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
)

var expectedConfig = Config{
	Host: "http://localhost",
	Port: 8080,
	Flag: true,
}

type ParameterStoreClientMock struct{}

type ParameterStoreUnauthorizedClientMock struct{}

func (ParameterStoreClientMock) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	value, _ := json.Marshal(&expectedConfig)
	stringValue := string(value)

	return &ssm.GetParameterOutput{
		Parameter: &types.Parameter{
			ARN:              new(string),
			DataType:         new(string),
			LastModifiedDate: &time.Time{},
			Name:             new(string),
			Selector:         new(string),
			SourceResult:     new(string),
			Type:             "",
			Value:            &stringValue,
			Version:          0,
		},
		ResultMetadata: middleware.Metadata{},
	}, nil
}

func (ParameterStoreUnauthorizedClientMock) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	return nil, errors.New("operation error SSM: GetParameter, failed to sign request: failed to retrieve credentials")
}

type Config struct {
	Host string `json:"host"`
	Port int    `json:"times"`
	Flag bool   `json:"flag"`
}

func TestReadConfigurationNoOverride(t *testing.T) {
	ctx := context.Background()
	mock := ParameterStoreClientMock{}

	config, err := ReadConfiguration[Config](ctx, mock, "config")
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, expectedConfig, *config)
}

func TestReadConfigurationWithOverride(t *testing.T) {
	ctx := context.Background()
	mock := ParameterStoreClientMock{}

	expectedOverrideConfig := Config{
		Host: "http://google.no",
		Port: 8080,
		Flag: false,
	}

	os.Setenv("host", "http://google.no")
	os.Setenv("flag", "false")

	config, err := ReadConfiguration[Config](ctx, mock, "config")
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, expectedOverrideConfig, *config)
}

func TestReadConfigurationUsesOverrideIfUnauthorized(t *testing.T) {
	ctx := context.Background()
	mock := ParameterStoreUnauthorizedClientMock{}

	expectedConfig := Config{
		Host: "http://google.no",
		Port: 8080,
		Flag: false,
	}

	os.Setenv("host", "http://google.no")
	os.Setenv("times", "8080")
	os.Setenv("flag", "false")

	config, err := ReadConfiguration[Config](ctx, mock, "config")
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, expectedConfig, *config)
}

func TestPanicIfConfigFieldMissing(t *testing.T) {
	ctx := context.Background()
	mock := ParameterStoreUnauthorizedClientMock{}

	os.Setenv("host", "http://google.no")
	os.Setenv("times", "8080")
	os.Unsetenv("flag")

	assert.Panics(t, func() { ReadConfiguration[Config](ctx, mock, "config") })
}
