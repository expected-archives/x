package xweb

import (
	"context"
	"errors"
	"github.com/caumette-co/x/xfoundation/contracts"
	"net/http"
)

type RenderViewBuilder struct {
	data   any
	engine string
	name   string
	layout *string

	*baseRenderBuilder[*RenderViewBuilder]
}

func ViewRender(name string) *RenderViewBuilder {
	renderViewBuilder := &RenderViewBuilder{
		engine: "x.gohtml",
		name:   name,
	}

	renderViewBuilder.baseRenderBuilder = newBaseRenderBuilder[*RenderViewBuilder](renderViewBuilder)

	return renderViewBuilder.
		WithContentType("text/html; charset=utf-8").
		WithStatusCode(http.StatusOK)
}

type ctxKeyEngines struct{}

var ctxKeyEnginesValue = ctxKeyEngines{}

func (r *RenderViewBuilder) WithData(data any) *RenderViewBuilder {
	r.data = data
	return r
}

func (r *RenderViewBuilder) WithEngine(engine string) *RenderViewBuilder {
	r.engine = engine
	return r
}

func (r *RenderViewBuilder) WithLayout(layout string) *RenderViewBuilder {
	r.layout = &layout
	return r
}

var ErrNoTemplateEngine = errors.New("no template engine")

type RenderData struct {
	Data   any
	Errors *contracts.ValidationError
}

func (r *RenderViewBuilder) Write(ctx context.Context, w http.ResponseWriter) error {
	m, ok := ctx.Value(ctxKeyEnginesValue).(map[string]contracts.TemplateEngine)
	if !ok {
		return ErrNoTemplateEngine
	}

	engine, ok := m[r.engine]
	if !ok {
		return ErrNoTemplateEngine
	}

	opts := make([]contracts.TemplateEngineOptsApplier, 0)
	if r.layout != nil {
		opts = append(opts, contracts.WithTemplateEngineOptLayout(*r.layout))
	}

	opts = append(opts, contracts.WithTemplateEngineOptData(RenderData{
		Data:   r.data,
		Errors: r.validationErrors,
	}))

	r.baseRenderBuilder.write(w)

	return engine.Execute(w, r.name, opts...)
}
