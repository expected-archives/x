package old

type errorWithResponse struct {
	error
	response Response
}

func WrapErrorWithResponse(err error, res Response) error {
	return &errorWithResponse{error: err, response: res}
}
