package xweb

import (
	"context"
	"encoding/json"
	"net/http"
)

type RenderJSONBuilder struct {
	data any

	*baseRenderBuilder[*RenderJSONBuilder]
}

func JSONRender() *RenderJSONBuilder {
	jsonBuilder := &RenderJSONBuilder{}
	jsonBuilder.baseRenderBuilder = newBaseRenderBuilder[*RenderJSONBuilder](jsonBuilder)

	return jsonBuilder.
		WithContentType("application/json; charset=utf-8").
		WithStatusCode(http.StatusOK)
}

func (r *RenderJSONBuilder) WithData(data any) *RenderJSONBuilder {
	r.data = data
	return r
}

func (r RenderJSONBuilder) Write(_ context.Context, w http.ResponseWriter) error {
	response := make(map[string]any)

	if r.validationErrors != nil {
		response["errors"] = r.validationErrors
	}

	if r.data != nil {
		response["data"] = r.data
	}

	r.baseRenderBuilder.write(w)

	return json.NewEncoder(w).Encode(r.data)
}
