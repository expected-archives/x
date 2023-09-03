package xweb

import (
	"github.com/caumette-co/x/xweb/binding"
	"net/http"
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
	binder := binding.NewBinder(binding.StringsParamExtractors, binding.ValuesParamExtractors) // TODO - move it :)
	// todo: ALEXIS :o /!\ URGENT PRIO POUR HIER
	_ = binder.Bind(r.Request, params)
	return params, nil
}
