package go_simple_router

type route struct {
	method  string
	pattern string
	handler handler
}
