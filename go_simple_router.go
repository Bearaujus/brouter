package go_simple_router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func NewGoSimpleRouter() Router {
	var chiRouter = chi.NewRouter()
	chiRouter.MethodNotAllowed(http.NotFound)
	chiRouter.Use(middleware.Logger)

	return &router{
		router:             chiRouter,
		routes:             make([]route, 0),
		fileServerPatterns: make([]string, 0),
	}
}

func NewGoSimpleRouterWithParam(chiRouter *chi.Mux) Router {
	if chiRouter == nil {
		return NewGoSimpleRouter()
	}

	return &router{
		router:             chiRouter,
		routes:             make([]route, 0),
		fileServerPatterns: make([]string, 0),
	}
}
