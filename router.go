package go_simple_router

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jedib0t/go-pretty/v6/table"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Router interface {
	Route(structRoute StructRoute)
	Routes(structRoutes []StructRoute)
	RouteFileServer(structRouteFileServer StructRouteFileServer) error
	RoutesFileServer(structRoutesFileServer []StructRouteFileServer) error
	Serve(host string, port int) error
}

type router struct {
	router             *chi.Mux
	routes             []route
	fileServerPatterns []string
}

func (r *router) Route(structRoute StructRoute) {
	// for each method
	for _, method := range structRoute.Methods {
		// setup handler
		var h handler
		if structRoute.HandlerFunc != nil {
			h.Func = structRoute.HandlerFunc
		}
		if structRoute.HandlerErrorFunc != nil {
			h.ErrorFunc = structRoute.HandlerErrorFunc
		}
		if structRoute.HandlerSuccessFunc != nil {
			h.SuccessFunc = structRoute.HandlerSuccessFunc
		}

		// add route data
		r.routes = append(r.routes, route{
			pattern: structRoute.Pattern,
			method:  method,
			handler: h,
		})
	}
}

func (r *router) Routes(structRoutes []StructRoute) {
	// for each struct route
	for _, structRoute := range structRoutes {
		r.Route(structRoute)
	}
}

func (r *router) RouteFileServer(structRouteFileServer StructRouteFileServer) error {
	// validate pattern
	if strings.ContainsAny(structRouteFileServer.Pattern, "{}*") {
		return errors.New("file server does not permit any url parameters")
	}

	// map pattern
	if structRouteFileServer.Pattern != "/" && structRouteFileServer.Pattern[len(structRouteFileServer.Pattern)-1] != '/' {
		r.router.Get(structRouteFileServer.Pattern, http.RedirectHandler(structRouteFileServer.Pattern+"/", http.StatusTemporaryRedirect).ServeHTTP)
		structRouteFileServer.Pattern += "/"
	}
	rawPattern := structRouteFileServer.Pattern
	structRouteFileServer.Pattern += "*"

	fsHandler := func(w http.ResponseWriter, r *http.Request) {
		pathPrefix := strings.TrimSuffix(chi.RouteContext(r.Context()).RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(structRouteFileServer.DirPath)))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}

	// add files pattern to router
	r.router.Get(structRouteFileServer.Pattern, fsHandler)
	r.fileServerPatterns = append(r.fileServerPatterns, rawPattern)

	return nil
}

func (r *router) RoutesFileServer(structRoutesFileServer []StructRouteFileServer) error {
	// for each struct route file server
	for _, structRouteFileServer := range structRoutesFileServer {
		if err := r.RouteFileServer(structRouteFileServer); err != nil {
			return err
		}
	}

	return nil
}

func (r *router) Serve(host string, port int) error {
	// init table writer
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ROUTE", "METHOD"})

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
		sort.Slice(v, func(i, j int) bool {
			return v[i] < v[j]
		})
		t.AppendRow(table.Row{k, v}, table.RowConfig{})
	}

	// log file server
	for _, fsPattern := range r.fileServerPatterns {
		t.AppendRow(table.Row{fsPattern, []string{http.MethodGet}}, table.RowConfig{})
	}

	// serve http
	t.SortBy([]table.SortBy{{Number: 1, Mode: table.Asc}})
	t.Render()
	fmt.Printf("> service running at %v\n", address)
	if err := http.Serve(listener, r.router); err != nil {
		return err
	}

	return nil
}
