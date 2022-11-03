package brouter

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func NewBRouter() Router {
	var chiRouter = chi.NewRouter()
	chiRouter.MethodNotAllowed(http.NotFound)
	chiRouter.Use(middleware.Logger)

	return &router{
		router:             chiRouter,
		routes:             make([]route, 0),
		fileServerPatterns: make([]string, 0),
	}
}

func NewBRouterWithParam(chiRouter *chi.Mux) Router {
	if chiRouter == nil {
		return NewBRouter()
	}

	return &router{
		router:             chiRouter,
		routes:             make([]route, 0),
		fileServerPatterns: make([]string, 0),
	}
}
