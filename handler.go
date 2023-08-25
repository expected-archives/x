package dise

import (
	"fmt"
	"net/http"
)

type Handler func(r *Request) (Response, error)

func Wrap(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := h(&Request{r})
		if err != nil {
			if wrappedError, ok := err.(*errorWithResponse); ok {
				response = wrappedError.response
			} else {
				LogError(r.Context(), err)
				response = &JSONResponse{
					StatusCode: http.StatusInternalServerError,
					Payload:    ErrorPayload{Message: http.StatusText(http.StatusInternalServerError)},
				}
			}
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
			LogError(r.Context(), fmt.Errorf("response.WriteBody: %v", err))
		}
	}
}
