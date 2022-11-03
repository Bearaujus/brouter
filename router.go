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
	Route(sr *StructRoute)
	RouteFileServer(srfs *StructRouteFileServer) error
	Serve(host string, port int) error
}

type router struct {
	router             *chi.Mux
	routes             []route
	fileServerPatterns []string
}

func (r *router) Route(sr *StructRoute) {
	// validate route data
	if sr == nil {
		return
	}

	// for each method
	for _, method := range sr.Methods {
		// setup handler
		var h handler
		if sr.HandlerFunc != nil {
			h.Func = sr.HandlerFunc
		}
		if sr.HandlerErrorFunc != nil {
			h.ErrorFunc = sr.HandlerErrorFunc
		}
		if sr.HandlerSuccessFunc != nil {
			h.SuccessFunc = sr.HandlerSuccessFunc
		}

		// add route data
		r.routes = append(r.routes, route{
			pattern: sr.Pattern,
			method:  method,
			handler: h,
		})
	}
}

func (r *router) RouteFileServer(srfs *StructRouteFileServer) error {
	// validate route data
	if srfs == nil {
		return errors.New("nil route file server")
	}

	// validate pattern
	if strings.ContainsAny(srfs.Pattern, "{}*") {
		return errors.New("file server does not permit any url parameters")
	}

	// map pattern
	if srfs.Pattern != "/" && srfs.Pattern[len(srfs.Pattern)-1] != '/' {
		r.router.Get(srfs.Pattern, http.RedirectHandler(srfs.Pattern+"/", http.StatusTemporaryRedirect).ServeHTTP)
		srfs.Pattern += "/"
	}
	rawPattern := srfs.Pattern
	srfs.Pattern += "*"

	fsHandler := func(w http.ResponseWriter, r *http.Request) {
		pathPrefix := strings.TrimSuffix(chi.RouteContext(r.Context()).RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(srfs.DirPath)))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}

	// add files pattern to router
	r.router.Get(srfs.Pattern, fsHandler)
	r.fileServerPatterns = append(r.fileServerPatterns, rawPattern)

	return nil
}

func (r *router) Serve(host string, port int) error {
	// map routes
	routesMap := map[string][]string{}
	for i := range r.routes {
		selectedRoute := r.routes[i]
		if selectedRoute.handler.Func == nil {
			selectedRoute.handler.Func = defaultHandlerFunc
		}
		if selectedRoute.handler.ErrorFunc == nil {
			selectedRoute.handler.ErrorFunc = defaultHandlerErrorFunc
		}
		if selectedRoute.handler.SuccessFunc == nil {
			selectedRoute.handler.SuccessFunc = defaultHandlerSuccessFunc
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
	logrus.Info("-------------------------------------------------------")
	logrus.Infof("service running at http://%v\n", address)
	logrus.Info("-------------------------------------------------------")
	if err := http.Serve(listener, r.router); err != nil {
		return err
	}

	return nil
}
