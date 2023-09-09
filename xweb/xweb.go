package xweb

import (
	"context"
	"github.com/caumette-co/x/xfoundation"
	"github.com/caumette-co/x/xfoundation/contracts"
	"github.com/expectedsh/dig"
	"go.uber.org/zap"
	"net"
	"net/http"
)

type Web struct {
	Config
	renderer map[string]contracts.Renderer
	router   *Router
}

type Config struct {
	Addr    string
	Routes  func(router *Router)
	BaseURL string
}

type Params struct {
	dig.In
	Renderer []contracts.Renderer `group:"x.renderer"`
	App      *xfoundation.App
}

var _ xfoundation.AppHandler = (*Web)(nil)

// New returns a new Web
func New(config Config) func(params Params) (xfoundation.AppHandlerOut, error) {
	return func(params Params) (xfoundation.AppHandlerOut, error) {
		w := &Web{
			Config:   config,
			renderer: make(map[string]contracts.Renderer),
			router:   nil,
		}
		router := newRouter(params.App, w)
		w.router = router

		for _, renderer := range params.Renderer {
			w.renderer[renderer.Name()] = renderer
		}

		return xfoundation.NewAppHandlerOut(w), nil
	}
}

// OnStartup implements xfoundation.AppHandler
func (w *Web) OnStartup(app *xfoundation.App) error {
	httpServer := &http.Server{Handler: w.router.handler}

	app.OnStart(func(ctx context.Context) error {
		w.Config.Routes(w.router)

		listener, err := net.Listen("tcp", w.Config.Addr)
		if err != nil {
			return err
		}
		app.Logger.Info("starting http server", zap.String("addr", listener.Addr().String()))
		go httpServer.Serve(listener)
		return nil
	})

	app.OnStop(func(ctx context.Context) error {
		return httpServer.Shutdown(ctx)
	})

	return nil
}
