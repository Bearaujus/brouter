package go_simple_router

import "net/http"

type handler struct {
	Func        func(w http.ResponseWriter, r *http.Request) (interface{}, error)
	ErrorFunc   func(w http.ResponseWriter, r *http.Request, err error)
	SuccessFunc func(w http.ResponseWriter, r *http.Request, data interface{})
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := h.Func(w, r)
	if err != nil {
		h.ErrorFunc(w, r, err)
		return
	}
	h.SuccessFunc(w, r, data)
}
