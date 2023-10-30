package awsssm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var sessionToken = os.Getenv("AWS_SESSION_TOKEN")

type GetParmeterExtensionApi interface {
	GetParameter(ctx context.Context, name string, decrypt bool) (*string, error)
}

type ParameterStoreExtensionClient struct{}

func (p *ParameterStoreExtensionClient) GetParameter(ctx context.Context, name string, decrypt bool) (*string, error) {
	// Define the URL
	url := fmt.Sprintf("http://localhost:2773/systemsmanager/parameters/get?name=%s&withDecryption=%s", url.QueryEscape(name), strconv.FormatBool(decrypt))

	// Create an HTTP GET request
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

	value := container["Parameter"].(map[string]any)
	valueString := value["Value"].(string)

	return &valueString, nil
}

func ReadParameterStoreParameterFromExtension(ctx context.Context, name string, api GetParmeterExtensionApi, decrypt bool) (*string, error) {
	return api.GetParameter(ctx, name, decrypt)
}
