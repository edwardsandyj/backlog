package router

import (
	"net/http"
	"regexp"
)

type Path struct {
	Pattern    *regexp.Regexp
	Handler    http.Handler
	HTTPMethod string
}

type Router struct {
	paths []*Path
}

func (r *Router) HandlerFunc(path string, httpmethod string, f func(http.ResponseWriter, *http.Request)) {
	r.paths = append(r.paths, &Path{
		Pattern:    regexp.MustCompile(path + "$"),
		Handler:    http.HandlerFunc(f),
		HTTPMethod: httpmethod,
	})
}
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, routing := range r.paths {
		if routing.HTTPMethod == request.Method && routing.Pattern.MatchString(request.URL.Path) {
			routing.Handler.ServeHTTP(writer, request)
			return
		}
	}
	http.NotFound(writer, request)
}
