package ginruntime

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oslokommune/common-lib-go/httpcomm"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler_ReturnsIndicatedStatus_WhenApiError(t *testing.T) {
	engine := New(context.Background())
	engine.Use(ErrorHandler())
	engine.AddRoute(nil, "/", GET, nil, func(c *gin.Context) {
		c.Error(Unauthorized("not logged in"))
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	engine.ServerHttp(res, req)

	assert.Equal(t, 401, res.Code)
}

func TestErrorHandler_Returns404_When404HttpCommError(t *testing.T) {
	engine := New(context.Background())
	engine.Use(ErrorHandler())
	engine.AddRoute(nil, "/", GET, nil, func(c *gin.Context) {
		c.Error(&httpcomm.HTTPError{Body: "test", StatusCode: 404})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	engine.ServerHttp(res, req)

	assert.Equal(t, 404, res.Code)
}

func TestErrorHandler_Returns403_When403HttpCommError(t *testing.T) {
	engine := New(context.Background())
	engine.Use(ErrorHandler())
	engine.AddRoute(nil, "/", GET, nil, func(c *gin.Context) {
		c.Error(&httpcomm.HTTPError{Body: "test", StatusCode: 403})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	engine.ServerHttp(res, req)

	assert.Equal(t, 403, res.Code)
}

func TestErrorHandler_Returns424_WhenGeneralHttpCommError(t *testing.T) {
	engine := New(context.Background())
	engine.Use(ErrorHandler())
	engine.AddRoute(nil, "/", GET, nil, func(c *gin.Context) {
		c.Error(&httpcomm.HTTPError{Body: "test", StatusCode: 429})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	engine.ServerHttp(res, req)
	assert.Equal(t, 424, res.Code)
}

func TestErrorHandler_Returns500_WhenGeneralError(t *testing.T) {
	engine := New(context.Background())
	engine.Use(ErrorHandler())
	engine.AddRoute(nil, "/", GET, nil, func(c *gin.Context) {
		c.Error(errors.New("error"))
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	engine.ServerHttp(res, req)
	assert.Equal(t, 500, res.Code)
}
