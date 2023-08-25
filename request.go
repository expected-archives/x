package dise

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	*http.Request
}

func (r *Request) ParseBody(body any) error {
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		return WrapErrorWithResponse(err, JSONResponse{
			Payload:    ErrorPayload{Message: "Invalid body."},
			StatusCode: http.StatusBadRequest,
		})
	}
	return nil
}
