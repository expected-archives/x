package old

import (
	"net/http"
	"strconv"
	"strings"
)

var CorsAllMethods = []string{
	http.MethodOptions,
	http.MethodHead,
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

type CorsOptions struct {
	Origin           string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

func Cors(opts *CorsOptions) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				w.Header().Add("Access-Control-Allow-Origin", opts.Origin)
				w.Header().Add("Access-Control-Allow-Credentials", strconv.FormatBool(opts.AllowCredentials))
				w.Header().Add("Access-Control-Allow-Methods", strings.Join(opts.AllowMethods, ","))
				w.Header().Add("Access-Control-Allow-Headers", strings.Join(opts.AllowHeaders, ","))
				return
			}

			w.Header().Add("Access-Control-Allow-Origin", opts.Origin)
			h.ServeHTTP(w, r)
		})
	}
}
