package xweb

import (
	"context"
	"github.com/caumette-co/x/xweb/binding"
	"go.uber.org/zap"
	"net/http"
)

type Handler[P any] func(r *Request[P]) (Response, error)

func (h Handler[P]) ServeHTTP(web *Web) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := newRequest[P](r, binding.Default) // TODO: inject binder in some way to allow customize it binder
		response, err := h(request)
		if err != nil {
			// TODO: handle error properly

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		ctx := r.Context()

		// provide templateEngines if any
		if web.templateEngines != nil {
			ctx = context.WithValue(ctx, ctxKeyEnginesValue, web.templateEngines)
		}

		if err := response.Write(ctx, w); err != nil {
			// todo: based on the http accept content type header, we should return a json response or a html response
			// with the error message
			zap.L().Error("error while writing response", zap.Error(err))
			return
		}
	}
}
