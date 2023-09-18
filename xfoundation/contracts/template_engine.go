package contracts

import (
	"github.com/expectedsh/dig"
	"github.com/samber/lo"
	"io"
)

type TemplateEngineOptions struct {
	Data   any
	Layout *string
}

type TemplateEngineOptsApplier func(*TemplateEngineOptions)

func WithTemplateEngineOptLayout(layout string) TemplateEngineOptsApplier {
	return func(options *TemplateEngineOptions) {
		options.Layout = lo.ToPtr(layout)
	}
}

func WithTemplateEngineOptData(data any) TemplateEngineOptsApplier {
	return func(options *TemplateEngineOptions) {
		options.Data = data
	}
}

// TemplateEngine is used to render templates
// group name: renderer
type TemplateEngine interface {
	Execute(writer io.Writer, templateName string, options ...TemplateEngineOptsApplier) error
	Name() string
}

// TemplateEngineOut is used to inject a TemplateEngine
type TemplateEngineOut struct {
	dig.Out

	Renderer TemplateEngine `group:"x.template_engine"`
}
