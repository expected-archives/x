package xweb

import (
	"github.com/caumette-co/x/xfoundation/contracts"
	"net/http"
)

type baseRenderBuilder[T Response] struct {
	statusCode       int
	headers          http.Header
	validationErrors *contracts.ValidationError

	ret T
}

func (r *baseRenderBuilder[T]) WithStatusCode(statusCode int) T {
	r.statusCode = statusCode
	return r.ret
}

func (r *baseRenderBuilder[T]) WithHeaders(headers http.Header) T {
	if r.headers == nil {
		r.headers = http.Header{}
	}

	for key, value := range headers {
		r.headers[key] = value
	}
	return r.ret
}

// WithErrors set the validation errors and the status code to 400
func (r *baseRenderBuilder[T]) WithErrors(errors *contracts.ValidationError) T {
	r.validationErrors = errors
	r.WithStatusCode(http.StatusBadRequest)
	return r.ret
}

func (r *baseRenderBuilder[T]) WithHeader(key, value string) T {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Set(key, value)
	return r.ret
}

func (r *baseRenderBuilder[T]) WithContentType(contentType string) T {
	return r.WithHeader("Content-Type", contentType)
}

func (r *baseRenderBuilder[T]) write(w http.ResponseWriter) {
	if r.headers != nil {
		for key, value := range r.headers {
			w.Header()[key] = value
		}
	}

	w.WriteHeader(r.statusCode)
}

func newBaseRenderBuilder[T Response](ret T) *baseRenderBuilder[T] {
	return &baseRenderBuilder[T]{ret: ret}
}
