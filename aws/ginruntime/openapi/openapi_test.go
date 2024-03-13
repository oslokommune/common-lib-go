package openapi

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestOpenAPIGenerateSchema(t *testing.T) {
	openapi := New("App", "1.0.0", "Test", "http://localhost:8080")
	openapi.Add("GET", "/", Annotate(Request[any](), Response[any](404), FileResponse(200, "application/json")))
	schema, err := openapi.MarshalJSON()
	assert.NoError(t, err)

	log.Info().Msg(string(schema))
	expected := "{\"openapi\":\"3.1.0\",\"info\":{\"title\":\"App\",\"description\":\"Test\",\"version\":\"1.0.0\"},\"paths\":{\"/\":{\"get\":{\"operationId\":\"GET-/\",\"responses\":{\"200\":{\"description\":\"OK\",\"content\":{\"application/json\":{\"schema\":{\"format\":\"binary\",\"type\":\"string\"}}}},\"404\":{\"description\":\"Not Found\",\"content\":{\"application/json\":{\"schema\":{\"type\":\"string\"}}}}}}}}}"
	assert.JSONEq(t, expected, string(schema))
}
