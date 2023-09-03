package handler

import (
	"github.com/caumette-co/x/xweb"
	"net/http"
)

func HandleHome(web *xweb.Provider) func(http.ResponseWriter, *http.Request) {
	//web.AddTemplate("layouts/landing.html", "home/index.html")
	//web.AddValidator("unique-user-email", func() {})

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}

func HandleDirect(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func HandleNew(r *xweb.Request[any]) (xweb.Response, error) {
	return xweb.JSONResponse{
		StatusCode: http.StatusOK,
		Payload:    map[string]interface{}{"hello": true},
	}, nil
}

func HandleNew2() xweb.Handler[any] {
	return func(r *xweb.Request[any]) (xweb.Response, error) {
		return xweb.JSONResponse{
			StatusCode: http.StatusOK,
			Payload:    map[string]interface{}{"hello": true},
		}, nil
	}
}
