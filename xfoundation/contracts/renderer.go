package contracts

import (
	"github.com/expectedsh/dig"
	"io"
)

type RendererOpts struct {
	Layout string
	Data   any
}

type RendererOptsApplier func(*RendererOpts)

func WithRendererApplierLayout(layout string) RendererOptsApplier {
	return func(options *RendererOpts) {
		options.Layout = layout
	}
}

func WithRendererApplierData(data any) RendererOptsApplier {
	return func(options *RendererOpts) {
		options.Data = data
	}
}

// Renderer is used to render templates
// group name: renderer
type Renderer interface {
	Render(writer io.Writer, name string, options ...RendererOptsApplier) error
	Name() string
}

// RendererOut is used to inject a Renderer
type RendererOut struct {
	dig.Out

	Renderer Renderer `group:"x.renderer"`
}
