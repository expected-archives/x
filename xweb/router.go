package xweb

import (
	"github.com/caumette-co/x/xfoundation"
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	app     *xfoundation.App
	handler *mux.Router
}

func newRouter(app *xfoundation.App) *Router {
	return &Router{
		app:     app,
		handler: mux.NewRouter(),
	}
}

func (r *Router) Route(method, path string, handler any) {
	r.handler.Methods(method).Path(path).HandlerFunc(handler.(func() http.HandlerFunc)())
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
