package xweb

import (
	"encoding/json"
	"io"
	"net/http"
)

type (
	Response interface {
		WriteBody(w io.Writer) error

		GetStatusCode() int

		GetHeaders() http.Header
	}

	JSONResponse struct {
		Payload    interface{}
		StatusCode int
		Headers    http.Header
	}
)

var _ Response = (*JSONResponse)(nil)

func (r JSONResponse) WriteBody(w io.Writer) error {
	b, err := json.Marshal(r.Payload)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (r JSONResponse) GetStatusCode() int {
	return r.StatusCode
}

func (r JSONResponse) GetHeaders() http.Header {
	return r.Headers
}
