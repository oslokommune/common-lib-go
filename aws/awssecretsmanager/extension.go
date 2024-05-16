package awssecretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var sessionToken = os.Getenv("AWS_SESSION_TOKEN")

type GetSecreteExtensionApi interface {
	GetSecret(ctx context.Context, name string, decrypt bool) (*string, error)
}

type SecretsManagerExtensionClient struct{}

func (p *SecretsManagerExtensionClient) GetParameter(ctx context.Context, name string, version *string) (map[string]any, error) {
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
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Aws-Parameters-Secrets-Token", sessionToken)

	// Call endpoint
	response, err := client.Do(req)
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

	var container map[string]any
	err = json.Unmarshal(responseBody, &container)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func ReadSecretsManagerSecretFromExtension(ctx context.Context, name string, api GetSecreteExtensionApi, decrypt bool) (*string, error) {
	return api.GetSecret(ctx, name, decrypt)
}
