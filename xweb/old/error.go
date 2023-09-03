package old

import "github.com/caumette-co/x/xweb"

type errorWithResponse struct {
	error
	response xweb.Response
}

func WrapErrorWithResponse(err error, res xweb.Response) error {
	return &errorWithResponse{error: err, response: res}
}
