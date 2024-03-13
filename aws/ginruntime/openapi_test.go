package ginruntime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oslokommune/common-lib-go/aws/ginruntime/openapi"
	"github.com/stretchr/testify/assert"
)

func TestOpenAPIEndpoint(t *testing.T) {
	engine := New(context.Background(), WithOpenAPI("test", "1.0", "description", "/"))

	type request struct {
		Body  string `json:"body" description:"description"`
		Query string `form:"form"`
		Uri   string `path:"path"`
	}

	annotations := openapi.Annotate(openapi.Request[request](), openapi.FileResponse(200, "application/json"))
	engine.AddRoute(nil, "/:path", GET, annotations, func(c *gin.Context) {})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/openapi.json", nil)
	engine.ServerHttp(res, req)
	expected := `
{
  "openapi": "3.1.0",
  "info": { "title": "test", "description": "description", "version": "1.0" },
  "paths": {
    "/{path}": {
      "get": {
        "operationId": "GET-/{path}",
        "parameters": [
          { "name": "form", "in": "query", "schema": { "type": "string" } },
          {
            "name": "path",
            "in": "path",
            "required": true,
            "schema": { "type": "string" }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": { "format": "binary", "type": "string" }
              }
            }
          }
        }
      }
    }
  }
}
	`
	assert.JSONEq(t, expected, res.Body.String())
}
