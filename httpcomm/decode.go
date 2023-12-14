package httpcomm

import (
	"context"
	"encoding/json"
)

func Decode[T any](ctx context.Context, responseBody []byte) (*T, error) {
	var message T
	err := json.Unmarshal([]byte(responseBody), &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func DecodeValue[T any](ctx context.Context, responseBody []byte) (T, error) {
	var message T
	err := json.Unmarshal([]byte(responseBody), &message)

	return message, err
}
