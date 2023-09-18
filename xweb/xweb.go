package xweb

import (
	"context"
	"crypto/tls"
	"github.com/caumette-co/x/xfoundation"
	"github.com/caumette-co/x/xfoundation/contracts"
	"github.com/expectedsh/dig"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type Web struct {
	Config
	router          *Router
	templateEngines map[string]contracts.TemplateEngine
	name            string
	logger          *zap.Logger
}

type Config struct {
	BaseURL string
	Routes  func(router *Router)

	// Name is the name of the web component
	// Can be API, Dashboard, CMS, Commerce ...
	Name string

	// Addr optionally specifies the TCP address for the server to listen on,
	// Example: ":8080"
	// If empty, ":http" (port 80) is used.
	Addr string

	// Below is almost a copy of the http.Server struct

	// DisableGeneralOptionsHandler, if true, passes "OPTIONS *" requests to the Handler,
	// otherwise responds with 200 OK and Content-Length: 0.
	DisableGeneralOptionsHandler bool

	// TLSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS. Note that this value is
	// cloned by ServeTLS and ListenAndServeTLS, so it's not
	// possible to modify the configuration with methods like
	// tls.Config.SetSessionTicketKeys. To use
	// SetSessionTicketKeys, use Server.Serve with a TLS Listener
	// instead.
	TLSConfig *tls.Config

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body. If ReadHeaderTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	ReadHeaderTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	// A zero or negative value means there will be no timeout.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeout time.Duration

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int
}

type Params struct {
	dig.In
	App *xfoundation.App

	TemplateEngines []contracts.TemplateEngine `group:"x.template_engine"`
}

var _ xfoundation.AppHandler = (*Web)(nil)

// New returns a new Web
func New(config Config) func(params Params) (xfoundation.AppHandlerOut, error) {
	return func(params Params) (xfoundation.AppHandlerOut, error) {
		name := config.Name
		if name == "" {
			name = "web"
		}

		w := &Web{
			Config:          config,
			router:          nil,
			logger:          zap.L().With(zap.String("component", name)),
			name:            name,
			templateEngines: map[string]contracts.TemplateEngine{},
		}

		router := newRouter(
			params.App,
			w)
		w.router = router

		for _, te := range params.TemplateEngines {
			w.templateEngines[te.Name()] = te
		}

		return xfoundation.NewAppHandlerOut(w), nil
	}
}

// OnStartup implements xfoundation.AppHandler
func (w *Web) OnStartup(app *xfoundation.App) error {
	httpServer := w.buildHTPServerFromConfig()

	app.OnStart(func(ctx context.Context) error {
		w.Config.Routes(w.router)

		addr := w.Config.Addr
		if addr == "" {
			addr = ":http"
		}

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		w.logger.Info("starting http server", zap.String("addr", listener.Addr().String()))
		go httpServer.Serve(listener)
		return nil
	})

	app.OnStop(func(ctx context.Context) error {
		return httpServer.Shutdown(ctx)
	})

	return nil
}

func (w *Web) buildHTPServerFromConfig() *http.Server {
	srv := &http.Server{Handler: w.router.handler}

	if w.Config.TLSConfig != nil {
		srv.TLSConfig = w.Config.TLSConfig
	}

	if w.Config.ReadTimeout != 0 {
		srv.ReadTimeout = w.Config.ReadTimeout
	}

	if w.Config.ReadHeaderTimeout != 0 {
		srv.ReadHeaderTimeout = w.Config.ReadHeaderTimeout
	}

	if w.Config.WriteTimeout != 0 {
		srv.WriteTimeout = w.Config.WriteTimeout
	}

	if w.Config.IdleTimeout != 0 {
		srv.IdleTimeout = w.Config.IdleTimeout
	}

	if w.Config.MaxHeaderBytes != 0 {
		srv.MaxHeaderBytes = w.Config.MaxHeaderBytes
	}

	return srv
}
