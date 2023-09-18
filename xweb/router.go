package xweb

import (
	"fmt"
	"github.com/caumette-co/x/xfoundation"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"reflect"
)

type Router struct {
	logger  *zap.Logger
	handler *chi.Mux
	web     *Web
	app     *xfoundation.App
}

func newRouter(app *xfoundation.App, web *Web) *Router {
	return &Router{
		handler: chi.NewRouter(),
		web:     web,
		logger:  web.logger.With(zap.String("subcomponent", "router")),
		app:     app,
	}
}

func (r *Router) Route(method, path string, handler any) {
	httpHandler, err := r.getHTTPHandler(handler)
	if err != nil {
		r.logger.Error("invalid handler",
			zap.Error(err),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("handler", reflect.TypeOf(handler).String()),
		)
		return
	}
	r.handler.Method(method, path, httpHandler)
}

func (r *Router) Use(handler ...any) {
	//r.middlewares = append(r.middlewares, handler...)
}

func (r *Router) Get(path string, handler any) {
	r.Route(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler any) {
	r.Route(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler any) {
	r.Route(http.MethodPut, path, handler)
}

func (r *Router) Patch(path string, handler any) {
	r.Route(http.MethodPatch, path, handler)
}

func (r *Router) Delete(path string, handler any) {
	r.Route(http.MethodDelete, path, handler)
}

type HTTPHandler interface {
	ServeHTTP(provider *Web) http.HandlerFunc
}

func (r *Router) getHTTPHandler(handler any) (http.Handler, error) {
	if value, ok := handler.(func(http.ResponseWriter, *http.Request)); ok {
		handler = http.HandlerFunc(value)
		return handler.(http.Handler), nil
	} else if _, ok := handler.(HTTPHandler); !ok {
		handlerType := reflect.TypeOf(handler)
		if handlerType.Kind() != reflect.Func {
			return nil, fmt.Errorf("handler is not a function")
		}

		if handlerType.NumOut() == 1 && handlerType.Out(0).Kind() == reflect.Func {
			values, err := r.app.Invoke(handler)
			if err != nil {
				return nil, fmt.Errorf("invoke: %w", err)
			} else if len(values) != 1 {
				return nil, fmt.Errorf("invoke: expected at least 1 value, got %d", len(values))
			}
			handler = values[0].Interface()
			handlerType = reflect.TypeOf(handler)
		}
	}

	if value, ok := handler.(HTTPHandler); ok {
		return value.ServeHTTP(r.web), nil
	} else {
		return nil, fmt.Errorf("handler is not a a web.HTTPHandler")
	}
}
