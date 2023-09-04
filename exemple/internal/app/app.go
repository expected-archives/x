package app

import (
	"github.com/caumette-co/x/example/internal/app/handler"
	"github.com/caumette-co/x/xfoundation"
	"github.com/caumette-co/x/xweb"
	"os"
)

var Default = xfoundation.App{
	Env: os.Getenv("APP_ENV"),
	Providers: []xfoundation.Provider{
		&xweb.Provider{
			Addr:   os.Getenv("ADDR"),
			Routes: Routes,
		},
	},
}

func Routes(router *xweb.Router) {
	router.Get("/", handler.HandleHome)
	router.Get("/direct", handler.HandleDirect)
	router.Get("/new", xweb.Handler[any](handler.HandleNew))
	router.Get("/new2", handler.HandleNew2)
	router.Get("/contact", xweb.Handler[handler.Contact](handler.HandleContact))
}
