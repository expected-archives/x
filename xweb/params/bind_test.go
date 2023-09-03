package params_test

import (
	"context"
	"github.com/gavv/httpexpect/v2"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"velo/modules/web"
	"velo/modules/web/params"
)

func TestBinder_Bind(t *testing.T) {
	mux := chi.NewRouter()

	type testStruct struct {
		FromPath                     string   `path:"fromPath"`
		FromQuery                    string   `query:"fromQuery"`
		FromHeader                   string   `header:"fromHeader"`
		FromContext                  string   `ctx:"fromContext"`
		FromJson                     string   `json:"fromJson"`
		FromDefaultValue             string   `query:"123" default:"value"`
		FromDefaultValueWithExploder []string `query:"123" default:"value1,value2,value3" exploder:","`
	}

	binder := params.NewBinder(
		web.StringsParamExtractors,
		web.ValuesParamExtractors,
	)

	mux.HandleFunc("/{fromPath}", func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "fromContext", "value"))

		var test testStruct
		err := binder.Bind(r, w, &test)
		require.NoError(t, err)

		require.Equal(t, "value", test.FromPath)
		require.Equal(t, "value", test.FromQuery)
		require.Equal(t, "value", test.FromHeader)
		require.Equal(t, "value", test.FromContext)
		require.Equal(t, "value", test.FromJson)
		require.Equal(t, "value", test.FromDefaultValue)
		require.Equal(t, []string{"value1", "value2", "value3"}, test.FromDefaultValueWithExploder)
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	e := httpexpect.Default(t, server.URL)
	e.GET("/{fromPath}").
		WithJSON(map[string]string{
			"fromJson": "value",
		}).
		WithHeader("fromHeader", "value").
		WithQuery("fromQuery", "value").
		WithHeader("fromContext", "value").
		WithPath("fromPath", "value").
		Expect()

}
