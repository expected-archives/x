package xweb

import (
	"encoding/json"
	"errors"
	"github.com/caumette-co/x/xfoundation/contracts"
	"io"
	"net/http"
)

type (
	Response interface {
		writeBody(_ *Web, w io.Writer) error

		GetStatusCode() int

		GetHeaders() http.Header
	}

	JSONResponse struct {
		Payload    interface{}
		StatusCode int
		Headers    http.Header
	}

	ViewResponse struct {
		Name       string
		Data       any
		Headers    http.Header
		StatusCode int

		// Layout is optional
		Layout string

		// Renderer is optional
		// It is used to choose the renderer to use, default value is set in DefaultViewRenderer
		Renderer string
	}
)

// DefaultViewRenderer is the default renderer used by the ViewResponse
var DefaultViewRenderer = "x.gohtml"

var _ Response = (*JSONResponse)(nil)

func (r JSONResponse) writeBody(_ *Web, w io.Writer) error {
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
	if r.Headers == nil {
		r.Headers = http.Header{}
	}

	r.Headers.Set("Content-Type", "application/json; charset=utf-8")

	return r.Headers
}

var _ Response = (*ViewResponse)(nil)

var ErrNoRenderer = errors.New("there is no renderer registered")

func (r ViewResponse) writeBody(p *Web, w io.Writer) error {
	renderer := r.getRenderer(p)
	if renderer == nil {
		return ErrNoRenderer
	}

	opts := make([]contracts.RendererOptsApplier, 0)
	if r.Layout != "" {
		opts = append(opts, contracts.WithRendererApplierLayout(r.Layout))
	}

	if r.Data != nil {
		opts = append(opts, contracts.WithRendererApplierData(r.Data))
	}

	return renderer.Render(w, r.Name, opts...)
}

func (r ViewResponse) getRenderer(p *Web) contracts.Renderer {
	if r.Renderer != "" {
		return p.renderer[r.Renderer]
	}

	return p.renderer[DefaultViewRenderer]
}

func (r ViewResponse) GetStatusCode() int {
	return r.StatusCode
}

func (r ViewResponse) GetHeaders() http.Header {
	if r.Headers == nil {
		r.Headers = http.Header{}
	}

	r.Headers.Set("Content-Type", "text/html; charset=utf-8")
	return r.Headers
}
