package go_simple_router

import "net/http"

type handler struct {
	Func          func(w http.ResponseWriter, r *http.Request) (interface{}, error)
	ErrorWriter   func(w http.ResponseWriter, r *http.Request, err error)
	SuccessWriter func(w http.ResponseWriter, r *http.Request, data interface{})
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := h.Func(w, r)
	if err != nil {
		h.ErrorWriter(w, r, err)
		return
	}
	h.SuccessWriter(w, r, data)
}
