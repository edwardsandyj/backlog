package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var casesTestDynamicRouter = []struct {
	name          string
	routeMethod   string
	routePattern  string
	requestMethod string
	requestURL    string
	want          int
}{
	{
		name:          "return StatusOK if route matches HTTP handler",
		routePattern:  "/userstories",
		routeMethod:   http.MethodGet,
		requestURL:    "/userstories",
		requestMethod: http.MethodGet,
		want:          http.StatusOK,
	},
	{
		name:          "return Not Found if route does not match HTTP handler",
		routePattern:  "/userstories",
		routeMethod:   http.MethodGet,
		requestURL:    "/userstories",
		requestMethod: http.MethodPost,
		want:          http.StatusNotFound,
	},
	{
		name:          "return StatusOK if route matches HTTP handler and Regexp",
		routePattern:  `/userstories/\d`,
		routeMethod:   http.MethodGet,
		requestURL:    "/userstories/1",
		requestMethod: http.MethodGet,
		want:          http.StatusOK,
	},
	{
		name:          "return Not Found if route not found",
		routePattern:  `/userstories\d`,
		routeMethod:   http.MethodPost,
		requestURL:    "/userstories/a",
		requestMethod: http.MethodPost,
		want:          http.StatusNotFound,
	},
}

func TestDynamicRouter(t *testing.T) {
	t.Log("Dynamic Router")
	for _, testcase := range casesTestDynamicRouter {
		t.Logf(testcase.name)
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(testcase.requestMethod, testcase.requestURL, nil)
		r := Router{}
		r.HandlerFunc(testcase.routePattern, testcase.routeMethod, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.ServeHTTP(recorder, request)
		if recorder.Code != testcase.want {
			t.Errorf("Output --> Got %d want %d", recorder.Code, testcase.want)
		}
	}
}
