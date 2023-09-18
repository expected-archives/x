package handler

import (
	"github.com/caumette-co/x/xweb"
	"net/http"
)

//
//func HandleHome() func(http.ResponseWriter, *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		_, _ = w.Write([]byte("ok"))
//	}
//}
//
//func HandleDirect(w http.ResponseWriter, r *http.Request) {
//	w.WriteHeader(http.StatusOK)
//	_, _ = w.Write([]byte("ok"))
//}
//
//func HandleNew(r *xweb.Request[any]) (xweb.Response, error) {
//	return xweb.JSONResponse{
//		StatusCode: http.StatusOK,
//		Payload:    map[string]interface{}{"hello": true},
//	}, nil
//}

func HandleJSON(r *xweb.Request[any]) (xweb.Response, error) {
	return xweb.JSONRender().WithData(map[string]interface{}{"hello": true}).WithStatusCode(http.StatusOK), nil

}

func HandleView(r *xweb.Request[any]) (xweb.Response, error) {
	return xweb.ViewRender("index").
		WithData(Contact{Email: 3}).
		WithStatusCode(http.StatusOK), nil
}

type Contact struct {
	Email int `query:"email"`
}

//
//func (c Contact) Validate() error {
//	return errors.New("not implemented")
//}
//
//func HandleContact(r *xweb.Request[Contact]) (xweb.Response, error) {
//	fmt.Println(r.Params().Email)
//
//	err := r.Validate()
//	if err != nil {
//		return nil, err
//	}
//
//	return xweb.JSONResponse{
//		StatusCode: http.StatusOK,
//		Payload:    map[string]interface{}{"hello": true},
//	}, nil
//}
