package contracts

import "net/http"

type HttpErrorHandler interface {
	Handle(err error, w http.ResponseWriter)
}
