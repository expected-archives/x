package old

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Request[S any, T any] struct {
	*http.Request
	state  *S
	params *T
}

func (r *Request[S, T]) State() *T {
	return r.params
}

func (r *Request[S, T]) Params() *T {
	return r.params
}

func (r *Request[S, T]) parseParams() (err error) {
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		if err = r.ParseMultipartForm(2 << 20); err == nil {

		}
	}

	if err != nil {
		err = WrapErrorWithResponse(err, JSONResponse{
			Payload:    ErrorPayload{Message: "Invalid body."},
			StatusCode: http.StatusBadRequest,
		})
	}
	return nil
}

func (r *Request[S, T]) ParseJsonBody(body any) (err error) {
	return json.NewDecoder(r.Body).Decode(body)
}
