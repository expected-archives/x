package xweb

import (
	"context"
	"net/http"
)

type Response interface {
	Write(ctx context.Context, w http.ResponseWriter) error
}
