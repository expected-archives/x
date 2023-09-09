package app

import (
	"github.com/caumette-co/x/example/internal/app/handler"
	"github.com/caumette-co/x/xfoundation"
	"github.com/caumette-co/x/xrenderer"
	"github.com/caumette-co/x/xweb"
	"os"
)

var Default = xfoundation.App{
	Env: xfoundation.AppEnv(os.Getenv("APP_ENV")),
	Providers: []any{
		xrenderer.New(xrenderer.Config{
			Folder:         "views",
			LayoutsFolder:  "layouts",
			PartialsFolder: "partials",
		}),
		xweb.New(xweb.Config{
			Addr:    ":8080",
			Routes:  Routes,
			BaseURL: "http://localhost:8080",
		}),
	},
}

func Routes(router *xweb.Router) {
	router.Get("/", handler.HandleHome)
	router.Get("/direct", handler.HandleDirect)
	router.Get("/new", xweb.Handler[any](handler.HandleNew))
	router.Get("/new2", handler.HandleNew2)
	router.Get("/contact", xweb.Handler[handler.Contact](handler.HandleContact))
	router.Get("/view", xweb.Handler[any](handler.HandleView))
}
