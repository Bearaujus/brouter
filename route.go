package go_simple_router

import "net/http"

type StructRoute struct {
	Pattern     string
	Methods     []string
	HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

	// optional (you can fill both, or one of them)
	HandlerErrorFunc   func(w http.ResponseWriter, r *http.Request, err error)
	HandlerSuccessFunc func(w http.ResponseWriter, r *http.Request, data interface{})
}

type StructRouteFileServer struct {
	Pattern string
	DirPath string
}

type route struct {
	pattern string
	method  string
	handler handler
}
