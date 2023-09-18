package app

import (
	"github.com/caumette-co/x/example/internal/app/handler"
	"github.com/caumette-co/x/xfoundation"
	"github.com/caumette-co/x/xweb"
	"os"
)

var Default = xfoundation.App{
	Env: xfoundation.Environment(os.Getenv("APP_ENV")),
	Providers: []any{
		xweb.New(xweb.Config{
			Addr:    ":3002",
			Routes:  Routes,
			BaseURL: "http://localhost:8080",
		}),
	},
}

func Routes(router *xweb.Router) {
	//router.Get("/", handler.HandleHome)
	//router.Get("/direct", handler.HandleDirect)
	//router.Get("/new", xweb.Handler[any](handler.HandleNew))
	//router.Get("/contact", xweb.Handler[handler.Contact](handler.HandleContact))
	router.Get("/json", xweb.Handler[any](handler.HandleJSON))
	router.Get("/view", xweb.Handler[any](handler.HandleView))
}
