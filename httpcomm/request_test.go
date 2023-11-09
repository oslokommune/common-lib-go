package httpcomm

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type HttpMockClient struct {
	mock.Mock
}

func (c *HttpMockClient) Do(req *http.Request) (*http.Response, error) {
	retvals := c.Called(req)
	return retvals.Get(0).(*http.Response), retvals.Error(1)
}

func TestCall200ReturnsResponseWithoutHttpError(t *testing.T) {
	ctx := context.Background()
	httpClient := &HttpMockClient{}
	response := http.Response{StatusCode: 200, Status: "200 OK", Body: io.NopCloser(strings.NewReader("OK"))}
	httpClient.On("Do", mock.Anything).Return(&response, nil)

	request := HTTPRequest{}
	res, err := Call(ctx, httpClient, request)

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Nil(t, res.Error)
}

func TestCall404ReturnsResponseWithHttpError(t *testing.T) {
	ctx := context.Background()
	httpClient := &HttpMockClient{}
	response := http.Response{StatusCode: 404, Status: "404 NOT FOUND", Body: io.NopCloser(strings.NewReader("NOT FOUND"))}
	httpClient.On("Do", mock.Anything).Return(&response, nil)

	request := HTTPRequest{}
	res, err := Call(ctx, httpClient, request)

	assert.Nil(t, err)
	assert.Equal(t, 404, res.StatusCode)
	assert.NotNil(t, res.Error)
	assert.Equal(t, 404, res.Error.StatusCode)
}

func TestCallFailingHttpDoReturnsError(t *testing.T) {
	ctx := context.Background()
	httpClient := &HttpMockClient{}
	err := errors.New("socket: connection dropped")
	httpClient.On("Do", mock.Anything).Return((*http.Response)(nil), err)

	request := HTTPRequest{}
	res, err := Call(ctx, httpClient, request)

	assert.Nil(t, res)
	assert.Error(t, err, "socket: connection dropped")
}
