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
	assert.JSONEq(t, `{}`, string(schema))
}
