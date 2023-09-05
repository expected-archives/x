package xweb

import (
	"context"
	"fmt"
	"github.com/caumette-co/x/xfoundation"
	"go.uber.org/zap"
	"net"
	"net/http"
)

type Provider struct {
	Addr string

	Routes func(router *Router)
}

func (p *Provider) Register(app *xfoundation.App) error {
	router := newRouter(app)
	if err := app.Provide(xfoundation.ProvideSingleValueFunc(router)); err != nil {
		return fmt.Errorf("failed to provide *webx.Router: %w", err)
	}
	if err := app.Provide(xfoundation.ProvideSingleValueFunc(p)); err != nil {
		return fmt.Errorf("failed to provide *webx.Provider: %w", err)
	}

	httpServer := &http.Server{Handler: router.handler}

	app.OnStart(func(ctx context.Context) error {
		p.Routes(router)
		listener, err := net.Listen("tcp", p.Addr)
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
