package xweb

import (
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
	handlerType := reflect.TypeOf(handler)
	log := r.app.Logger.With(
		zap.String("method", method), zap.String("path", path), zap.Any("handler", handlerType))
	if handlerType.Kind() != reflect.Func {
		log.Error("route handler must be a func")
		return
	}

	if handlerType.NumOut() == 1 && handlerType.Out(0).Kind() == reflect.Func {
		values, err := r.app.Invoke(handler)
		if err != nil {
			log.Error("failed to invoke route constructor", zap.Error(err))
			return
		} else if len(values) != 1 {
			log.Error("an handler func must be returned", zap.Error(err))
			return
		}
		handler = values[0].Interface()
		handlerType = reflect.TypeOf(handler)
	}

	if value, ok := handler.(func(http.ResponseWriter, *http.Request)); ok {
		handler = http.HandlerFunc(value)
	}
	r.handler.Method(method, path, handler.(http.Handler))
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
