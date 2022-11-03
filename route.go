package brouter

import "net/http"

type (
	StructRoute struct {
		Pattern     string
		Methods     []string
		HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

		// optional - if nil, will use default error/success func
		HandlerErrorFunc   func(w http.ResponseWriter, r *http.Request, err error)
		HandlerSuccessFunc func(w http.ResponseWriter, r *http.Request, data interface{})
	}

	StructRouteFileServer struct {
		Pattern string
		DirPath string
	}

	route struct {
		pattern string
		method  string
		handler *handler
	}
)
