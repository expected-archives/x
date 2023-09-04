package xweb

import (
	"github.com/caumette-co/x/xweb/binding"
	"net/http"
)

type Request[P any] struct {
	*http.Request
	params *P

	binder       binding.Binder
	bindingError error
}

func newRequest[P any](request *http.Request, binder binding.Binder) *Request[P] {
	return &Request[P]{
		Request: request,
		binder:  binder,
	}
}

func (r *Request[P]) Params() *P {
	if r.params == nil {
		r.params, r.bindingError = r.parseParams()
	}
	return r.params
}

func (r *Request[P]) BindingErrors() error {
	return r.bindingError
}

func (r *Request[P]) parseParams() (*P, error) {
	params := new(P)
	err := r.binder.Bind(r.Request, params)
	return params, err
}
