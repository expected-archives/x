package xweb

import (
	"net/http"
)

type Handler[P any] func(r *Request[P]) (Response, error)

func WrapHandler[P any](h Handler[P]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		xr := &Request[P]{Request: r}
		response, err := h(xr)
		if err != nil {
			//if wrappedError, ok := err.(*errorWithResponse); ok {
			//	response = wrappedError.response
			//} else {
			//	LogError(r.Context(), err)
			//	response = &xweb.JSONResponse{
			//		StatusCode: http.StatusInternalServerError,
			//		Payload:    ErrorPayload{Message: http.StatusText(http.StatusInternalServerError)},
			//	}
			//}
		}

		if response == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for key, value := range response.GetHeaders() {
			w.Header()[key] = value
		}

		w.WriteHeader(response.GetStatusCode())
		if err := response.WriteBody(w); err != nil {
			//LogError(r.Context(), fmt.Errorf("response.WriteBody: %v", err))
		}
	}
}
