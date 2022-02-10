package router

import (
	"net/http"
	"regexp"
)

type Routing struct {
	Pattern    *regexp.Regexp
	Handler    http.Handler
	HTTPMethod string
}

type Router struct {
	routings []*Routing
}

func (r *Router) HandlerFunc(p string, httpmethod string, f func(http.ResponseWriter, *http.Request)) {
	r.routings = append(r.routings, &Routing{
		Pattern:    regexp.MustCompile(p + "$"),
		Handler:    http.HandlerFunc(f),
		HTTPMethod: httpmethod,
	})
}
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, routing := range r.routings {
		if routing.HTTPMethod == request.Method && routing.Pattern.MatchString(request.URL.Path) {
			routing.Handler.ServeHTTP(writer, request)
			return
		}
	}
	http.NotFound(writer, request)
}
