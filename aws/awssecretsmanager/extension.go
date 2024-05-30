package awssecretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var sessionToken = os.Getenv("AWS_SESSION_TOKEN")

type GetSecretFromExtensionApi interface {
	GetSecret(ctx context.Context, name string, version *string) (*SecretData, error)
}

type SecretData struct {
	ARN            string         `json:"ARN"`
	Name           string         `json:"Name"`
	VersionId      string         `json:"VersionId"`
	SecretString   *string        `json:"SecretString"`
	SecretBinary   []byte         `json:"SecretBinary"`
	VersionStages  []string       `json:"VersionStages"`
	CreatedDate    time.Time      `json:"CreatedDate"`
	ResultMetadata map[string]any `json:"ResultMetadata"`
}

type SecretsManagerExtensionClient struct {
	httpClient *http.Client
	tracing    bool
}

var _ GetSecretFromExtensionApi = (*SecretsManagerExtensionClient)(nil)

func NewExtensionClient(tracing bool) *SecretsManagerExtensionClient {
	var httpClient *http.Client
	if tracing {
		commonLabels := []attribute.KeyValue{
			attribute.String("otel.resource.service.name", "secret manager extentsion client"),
		}

		httpClient = &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport, otelhttp.WithSpanOptions(trace.WithAttributes(commonLabels...))),
			Timeout:   1 * time.Second,
		}
	} else {
		httpClient = &http.Client{}
	}

	return &SecretsManagerExtensionClient{
		httpClient: httpClient,
		tracing:    tracing,
	}
}

func (p *SecretsManagerExtensionClient) GetSecret(ctx context.Context, name string, version *string) (*SecretData, error) {
	// Define the URL
	parsedURL, err := url.Parse("http://localhost:2773/secretsmanager/get")
	if err != nil {
		return nil, err
	}
	query := parsedURL.Query()
	query.Set("secretId", name)
	if version != nil {
		query.Set("versionId", *version)
	}
	parsedURL.RawQuery = query.Encode()

	// Create an HTTP GET request
	req, err := http.NewRequestWithContext(ctx, "GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Aws-Parameters-Secrets-Token", sessionToken)

	// Call endpoint
	response, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check if the response status code is not OK (200)
	if response.StatusCode != http.StatusOK {
		text, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("HTTP request failed with status code %d. Text: %s", response.StatusCode, text)
	}

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var container SecretData
	err = json.Unmarshal(responseBody, &container)
	if err != nil {
		return nil, err
	}

	return &container, nil
}

func ReadSecretsManagerSecretFromExtension(ctx context.Context, name string, api GetSecretFromExtensionApi, version *string) (*SecretData, error) {
	return api.GetSecret(ctx, name, version)
}
