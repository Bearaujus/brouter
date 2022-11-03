package brouter

import "net/http"

type handler struct {
	Func        func(w http.ResponseWriter, r *http.Request) (interface{}, error)
	ErrorFunc   func(w http.ResponseWriter, r *http.Request, err error)
	SuccessFunc func(w http.ResponseWriter, r *http.Request, data interface{})
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := h.Func(w, r)
	if err != nil {
		h.ErrorFunc(w, r, err)
		return
	}
	h.SuccessFunc(w, r, data)
}

func (h *handler) route(structRoute StructRoute) {
	if structRoute.HandlerFunc != nil {
		h.Func = structRoute.HandlerFunc
	}

	if structRoute.HandlerErrorFunc != nil {
		h.ErrorFunc = structRoute.HandlerErrorFunc
	}

	if structRoute.HandlerSuccessFunc != nil {
		h.SuccessFunc = structRoute.HandlerSuccessFunc
	}
}

func (h *handler) setDefaultFunc(r *bRouter) {
	if h.Func == nil {
		h.Func = defaultHandlerFunc
	}

	if h.ErrorFunc == nil {
		h.ErrorFunc = defaultHandlerErrorFunc
		if r.defaultHandlerErrorFunc != nil {
			h.ErrorFunc = r.defaultHandlerErrorFunc
		}
	}

	if h.SuccessFunc == nil {
		h.SuccessFunc = defaultHandlerSuccessFunc
		if r.defaultHandlerSuccessFunc != nil {
			h.SuccessFunc = r.defaultHandlerSuccessFunc
		}
	}
}
