package params

import (
	"net/http"
)

// StringsParamExtractor extract multiples strings values from request.
type StringsParamExtractor interface {
	Extract(req *http.Request, valueOfTag string) ([]string, error)
	Tag() string
}

// ValueParamExtractor extract one value (a type) from http request.
type ValueParamExtractor interface {
	Extract(req *http.Request, valueOfTag string) (interface{}, error)
	Tag() string
}
