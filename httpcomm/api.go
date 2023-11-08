package httpComm

import (
	"context"
	"io"
)

type Api interface {
	Delete(ctx context.Context, resourceURL string, token *string) ([]byte, error)
	Get(ctx context.Context, resourceURL string, token *string) ([]byte, error)
	Put(ctx context.Context, resourceURL string, body io.Reader, contentType string, token *string) ([]byte, error)
	Post(ctx context.Context, resourceURL string, body io.Reader, contentType string, token *string) ([]byte, error)
}
