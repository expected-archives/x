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
	app     *xfoundation.App
	handler *chi.Mux
}

func newRouter(app *xfoundation.App) *Router {
	return &Router{
		app:     app,
		handler: chi.NewRouter(),
	}
}

func (r *Router) Route(method, path string, handler any) {
	httpHandler, err := r.getHttpHandler(handler)
	if err != nil {
		r.app.Logger.Error("invalid handler",
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

func (r *Router) getHttpHandler(handler any) (http.Handler, error) {
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

	if value, ok := handler.(func(http.ResponseWriter, *http.Request)); ok {
		handler = http.HandlerFunc(value)
	}

	value, ok := handler.(http.Handler)
	if !ok {
		return nil, fmt.Errorf("handler is not a http.Handler")
	}
	return value, nil
}
