package brouter

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jedib0t/go-pretty/v6/table"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
)

func NewBRouter() BRouter {
	var chiRouter = chi.NewRouter()
	chiRouter.MethodNotAllowed(http.NotFound)
	chiRouter.Use(middleware.Logger)

	return &bRouter{
		router:             chiRouter,
		routes:             make([]route, 0),
		fileServerPatterns: make([]string, 0),
	}
}

func NewBRouterWithParam(chiRouter *chi.Mux) BRouter {
	if chiRouter == nil {
		return NewBRouter()
	}

	return &bRouter{
		router:             chiRouter,
		routes:             make([]route, 0),
		fileServerPatterns: make([]string, 0),
	}
}

type (
	BRouter interface {
		Route(structRoute StructRoute)
		Routes(structRoutes []StructRoute)
		RouteFileServer(structRouteFileServer StructRouteFileServer) error
		RoutesFileServer(structRoutesFileServer []StructRouteFileServer) error
		Serve(host string, port int) error

		SetDefaultHandlerErrorFunc(defaultHandlerErrorFunc func(w http.ResponseWriter, r *http.Request, err error))
		SetDefaultHandlerSuccessFunc(defaultHandlerSuccessFunc func(w http.ResponseWriter, r *http.Request, data interface{}))
	}

	bRouter struct {
		router                    *chi.Mux
		routes                    []route
		fileServerPatterns        []string
		defaultHandlerErrorFunc   func(w http.ResponseWriter, r *http.Request, err error)
		defaultHandlerSuccessFunc func(w http.ResponseWriter, r *http.Request, data interface{})
	}
)

func (r *bRouter) Route(structRoute StructRoute) {
	// for each method
	for _, method := range structRoute.Methods {
		// setup handler
		var h = &handler{}
		h.route(structRoute)

		// add route data
		r.routes = append(r.routes, route{
			pattern: structRoute.Pattern,
			method:  method,
			handler: h,
		})
	}
}

func (r *bRouter) Routes(structRoutes []StructRoute) {
	// for each struct route
	for _, structRoute := range structRoutes {
		r.Route(structRoute)
	}
}

func (r *bRouter) RouteFileServer(structRouteFileServer StructRouteFileServer) error {
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

	// add files pattern to bRouter
	r.router.Get(structRouteFileServer.Pattern, fsHandler)
	r.fileServerPatterns = append(r.fileServerPatterns, rawPattern)

	return nil
}

func (r *bRouter) RoutesFileServer(structRoutesFileServer []StructRouteFileServer) error {
	// for each struct route file server
	for _, structRouteFileServer := range structRoutesFileServer {
		if err := r.RouteFileServer(structRouteFileServer); err != nil {
			return err
		}
	}

	return nil
}

func (r *bRouter) Serve(host string, port int) error {
	// init table writer
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ROUTE", "METHOD", "TYPE"})

	// map routes
	routesMap := map[string][]string{}
	for i := range r.routes {
		selectedRoute := r.routes[i]
		selectedRoute.handler.setDefaultFunc(r)
		r.router.Method(selectedRoute.method, selectedRoute.pattern, selectedRoute.handler)
		routesMap[selectedRoute.pattern] = append(routesMap[selectedRoute.pattern], selectedRoute.method)
	}

	// map host address
	address := fmt.Sprintf("%v:%v", host, port)
	listener, errListener := net.Listen("tcp", address)
	if errListener != nil {
		return errListener
	}

	// log routing tables
	for k, v := range routesMap {
		sort.Slice(v, func(i, j int) bool {
			return v[i] < v[j]
		})
		t.AppendRow(table.Row{k, v, "basic"}, table.RowConfig{})
	}

	for _, fsPattern := range r.fileServerPatterns {
		t.AppendRow(table.Row{fsPattern, []string{http.MethodGet}, "fs"}, table.RowConfig{})
	}

	t.SortBy([]table.SortBy{{Number: 3, Mode: table.Asc}, {Number: 1, Mode: table.Asc}})
	t.Render()

	// serve http
	fmt.Printf("service running at %v\n", address)
	if err := http.Serve(listener, r.router); err != nil {
		return err
	}

	return nil
}

func (r *bRouter) SetDefaultHandlerErrorFunc(defaultHandlerErrorFunc func(w http.ResponseWriter, r *http.Request, err error)) {
	r.defaultHandlerErrorFunc = defaultHandlerErrorFunc
}
func (r *bRouter) SetDefaultHandlerSuccessFunc(defaultHandlerSuccessFunc func(w http.ResponseWriter, r *http.Request, data interface{})) {
	r.defaultHandlerSuccessFunc = defaultHandlerSuccessFunc
}
