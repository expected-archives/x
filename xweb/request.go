package xweb

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Request[P any] struct {
	*http.Request
	params      *P
	paramsError error
}

func (r *Request[P]) Params() *P {
	if r.params == nil {
		r.params, r.paramsError = r.parseParams()
	}
	return r.params
}

func (r *Request[P]) Valid() error {
	return r.paramsError
}

func (r *Request[P]) parseParams() (*P, error) {
	params := new(P)
	if contentType := r.Header.Get("Content-Type"); strings.HasPrefix(contentType, "application/json") {
		return params, json.NewDecoder(r.Body).Decode(params)
	}
	return params, nil
}
