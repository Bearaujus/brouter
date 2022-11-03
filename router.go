package go_simple_router

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
)

type Router interface {
	Route(pattern string, methods []string, handlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error), errorHandler func(w http.ResponseWriter, r *http.Request, err error), successHandler func(w http.ResponseWriter, r *http.Request, data interface{}))
	RouteFileServer(pattern string, fileServerDirPath string) error
	Serve(host string, port int) error
}

type router struct {
	router             *chi.Mux
	routes             []route
	fileServerPatterns []string
}

func (r *router) Route(pattern string, methods []string, handlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error), errorWriter func(w http.ResponseWriter, r *http.Request, err error), successWriter func(w http.ResponseWriter, r *http.Request, data interface{})) {
	// for each method
	for _, method := range methods {
		// setup handler
		h := handler{
			Func: handlerFunc,
		}
		if errorWriter != nil {
			h.ErrorWriter = errorWriter
		}
		if successWriter != nil {
			h.SuccessWriter = successWriter
		}

		// add route data
		r.routes = append(r.routes, route{
			method:  method,
			pattern: pattern,
			handler: h,
		})
	}
}

func (r *router) RouteFileServer(pattern string, fileServerDirPath string) error {
	// validate pattern
	if strings.ContainsAny(pattern, "{}*") {
		return errors.New("file server does not permit any url parameters")
	}

	// map pattern
	if pattern != "/" && pattern[len(pattern)-1] != '/' {
		r.router.Get(pattern, http.RedirectHandler(pattern+"/", http.StatusTemporaryRedirect).ServeHTTP)
		pattern += "/"
	}
	rawPattern := pattern
	pattern += "*"

	fsHandler := func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(fileServerDirPath)))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}

	// add files pattern to router
	r.router.Get(pattern, fsHandler)
	r.fileServerPatterns = append(r.fileServerPatterns, rawPattern)

	return nil
}

func (r *router) Serve(host string, port int) error {
	// map routes
	routesMap := map[string][]string{}
	for i := range r.routes {
		selectedRoute := r.routes[i]
		if selectedRoute.handler.ErrorWriter == nil {
			selectedRoute.handler.ErrorWriter = DefaultErrorWriter
		}
		if selectedRoute.handler.SuccessWriter == nil {
			selectedRoute.handler.SuccessWriter = DefaultSuccessWriter
		}
		r.router.Method(selectedRoute.method, selectedRoute.pattern, selectedRoute.handler)
		routesMap[selectedRoute.pattern] = append(routesMap[selectedRoute.pattern], selectedRoute.method)
	}

	// map host address
	address := fmt.Sprintf("%v:%v", host, port)
	listener, errListener := net.Listen("tcp", address)
	if errListener != nil {
		return errListener
	}

	// log routes
	for k, v := range routesMap {
		logrus.Infof("route: %v http://%v%v", v, address, k)
	}

	// log file server
	for _, fsPattern := range r.fileServerPatterns {
		logrus.Infof("route-fs: http://%v%v\n", address, fsPattern)
	}

	// serve http
	logrus.Info("-----------------------------------------")
	logrus.Infof("service running at http://%v\n", address)
	if err := http.Serve(listener, r.router); err != nil {
		return err
	}

	return nil
}
