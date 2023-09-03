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
