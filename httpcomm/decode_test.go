package httpcomm

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Field1 string
	Field2 int
}

func TestDecodeValidJsonToStringReturnsValue(t *testing.T) {
	actual := "Hello world"
	bytes, _ := json.Marshal(actual)

	received, err := Decode[string](context.Background(), bytes, false)

	assert.Nil(t, err)
	assert.Equal(t, actual, *received)
}

func TestDecodeInvalidJsonToStringReturnsError(t *testing.T) {
	ctx := context.Background()

	invalid := 10
	bytes, _ := json.Marshal(invalid)

	received, err := Decode[string](ctx, bytes, false)

	assert.Nil(t, received)
	assert.Error(t, err)
}

func TestDecodeValidJsonToStructReturnsValue(t *testing.T) {
	actual := TestStruct{"Hello world", 10}
	bytes, _ := json.Marshal(actual)

	received, err := Decode[TestStruct](context.Background(), bytes, false)

	assert.Nil(t, err)
	assert.Equal(t, actual, *received)
}

func TestDecodeInvalidJsonToStructReturnsValue(t *testing.T) {
	invalid := 10
	bytes, _ := json.Marshal(invalid)

	received, err := Decode[TestStruct](context.Background(), bytes, false)

	assert.Nil(t, received)
	assert.Error(t, err)
}
