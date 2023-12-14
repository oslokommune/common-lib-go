package httpcomm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/rs/zerolog/log"
)

// Interface for performing HTTP calls, e.g. a `http.Client`
type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

func CreateRequest(ctx context.Context, httpRequest HTTPRequest) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, httpRequest.Method, httpRequest.Url, httpRequest.Body)
	if err != nil {
		return nil, err
	}

	for key, value := range httpRequest.Headers {
		req.Header.Set(key, value)
	}

	if httpRequest.Token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *httpRequest.Token))
	}

	return req, nil
}

func Call(ctx context.Context, httpClient HttpDoer, httpRequest HTTPRequest) (*HTTPResponse, error) {
	req, err := CreateRequest(ctx, httpRequest)
	if err != nil {
		return nil, err
	}

	// Logs request if debug level is enabled
	if log.Debug().Enabled() {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			log.Debug().Msg(string(reqDump))
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	// Logs response if trace level is enabled
	if log.Trace().Enabled() {
		resDump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			log.Trace().Msg(string(resDump))
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, &HTTPError{
			Body:       string(body),
			StatusCode: statusCode,
		}
	}

	return &HTTPResponse{
		StatusCode: statusCode,
		Body:       string(body),
	}, nil
}
