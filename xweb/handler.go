package xweb

import (
	"fmt"
	"github.com/caumette-co/x/xweb/binding"
	"net/http"
)

type Handler[P any] func(r *Request[P]) (Response, error)

func (h Handler[P]) ServeHTTP(provider *Web) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		xr := newRequest[P](r, binding.Default) // TODO: inject binder in some way to allow customize it binder
		response, err := h(xr)
		if err != nil {
			// todo:
			// - read accept content type
			//if wrappedError, ok := err.(*errorWithResponse); ok {
			//	response = wrappedError.response
			//} else {
			//	LogError(r.Context(), err)
			//	response = &xweb.JSONResponse{
			//		StatusCode: http.StatusInternalServerError,
			//		Payload:    ErrorPayload{Message: http.StatusText(http.StatusInternalServerError)},
			//	}
			//}

			// TODO: handle error properly

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		if response == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for key, value := range response.GetHeaders() {
			w.Header()[key] = value
		}

		w.WriteHeader(response.GetStatusCode())
		if err := response.writeBody(provider, w); err != nil {
			// todo log error
			fmt.Println(err)
			return
		}
	}
}
